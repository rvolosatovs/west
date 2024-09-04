use core::time::Duration;

use anyhow::Context as _;
use wasi_preview1_component_adapter_provider::{
    WASI_SNAPSHOT_PREVIEW1_ADAPTER_NAME, WASI_SNAPSHOT_PREVIEW1_REACTOR_ADAPTER,
};
use wasmtime::component::{Component, Linker, Resource, ResourceTable, Type, TypedFunc, Val};
use wasmtime::{Engine, Store};
use wasmtime_wasi::{WasiCtx, WasiCtxBuilder, WasiView};
use wasmtime_wasi_http::types::HostIncomingRequest;
use wasmtime_wasi_http::{WasiHttpCtx, WasiHttpView};

mod bindings {
    wasmtime::component::bindgen!({
        trappable_imports: true,
        path: "../../wit",
        with: {
            "wasi:http/types@0.2.0/fields": wasmtime_wasi_http::bindings::http::types::Fields,
            "wasi:http/types@0.2.0/future-incoming-response": wasmtime_wasi_http::bindings::http::types::FutureIncomingResponse,
            "wasi:http/types@0.2.0/incoming-request": wasmtime_wasi_http::bindings::http::types::IncomingRequest,
            "wasi:http/types@0.2.0/outgoing-request": wasmtime_wasi_http::bindings::http::types::OutgoingRequest,
            "wasi:http/types@0.2.0/response-outparam": wasmtime_wasi_http::bindings::http::types::ResponseOutparam,
        },
        world: "http-test-passthrough",
    });
}

pub const PASSTHROUGH: &[u8] = include_bytes!(concat!(env!("OUT_DIR"), "/west_passthrough.wasm"));

struct Ctx {
    wasi: WasiCtx,
    http: WasiHttpCtx,
    table: ResourceTable,
}

impl WasiView for Ctx {
    fn ctx(&mut self) -> &mut WasiCtx {
        &mut self.wasi
    }
    fn table(&mut self) -> &mut ResourceTable {
        &mut self.table
    }
}

impl WasiHttpView for Ctx {
    fn ctx(&mut self) -> &mut WasiHttpCtx {
        &mut self.http
    }
    fn table(&mut self) -> &mut ResourceTable {
        &mut self.table
    }
}

impl bindings::west::test::http_test::Host for Ctx {
    fn new_response_outparam(
        &mut self,
    ) -> wasmtime::Result<(
        Resource<wasmtime_wasi_http::types::HostResponseOutparam>,
        Resource<wasmtime_wasi_http::types::HostFutureIncomingResponse>,
    )> {
        let (res_tx, res_rx) = tokio::sync::oneshot::channel();
        let out = WasiHttpView::new_response_outparam(self, res_tx)
            .context("failed to construct `response-outparam`")?;
        let res = WasiHttpView::table(self)
            .push(
                wasmtime_wasi_http::types::HostFutureIncomingResponse::Pending(
                    wasmtime_wasi::runtime::spawn(async {
                        match res_rx.await.context("failed to receive response")? {
                            Ok(resp) => Ok(Ok(wasmtime_wasi_http::types::IncomingResponse {
                                resp,
                                worker: None,
                                between_bytes_timeout: Duration::from_secs(1),
                            })),
                            Err(err) => Ok(Err(err)),
                        }
                    }),
                ),
            )
            .context("failed to push `future-incoming-response` into resource table")?;
        Ok((out, res))
    }

    fn new_incoming_request(
        &mut self,
        req: Resource<wasmtime_wasi_http::types::HostOutgoingRequest>,
    ) -> wasmtime::Result<Resource<wasmtime_wasi_http::types::HostIncomingRequest>> {
        let wasmtime_wasi_http::types::HostOutgoingRequest {
            method,
            scheme,
            authority,
            path_with_query,
            headers,
            body,
        } = WasiHttpView::table(self)
            .delete(req)
            .context("failed to delete outgoing request")?;

        let uri = http::Uri::builder();
        let uri = match &scheme {
            None | Some(wasmtime_wasi_http::bindings::http::types::Scheme::Http) => {
                uri.scheme(http::uri::Scheme::HTTP)
            }
            Some(wasmtime_wasi_http::bindings::http::types::Scheme::Https) => {
                uri.scheme(http::uri::Scheme::HTTPS)
            }
            Some(wasmtime_wasi_http::bindings::http::types::Scheme::Other(scheme)) => {
                uri.scheme(scheme.as_str())
            }
        };
        let uri = if let Some(path_with_query) = path_with_query {
            uri.path_and_query(path_with_query)
        } else {
            uri.path_and_query("/")
        };
        let uri = if let Some(authority) = authority {
            uri.authority(authority)
        } else {
            uri.authority("west")
        };
        let uri = uri.build().context("failed to build URI")?;
        let mut req = http::Request::builder();
        if let Some(h) = req.headers_mut() {
            *h = headers;
        }
        let req = match &method {
            wasmtime_wasi_http::bindings::http::types::Method::Get => req.method(http::Method::GET),
            wasmtime_wasi_http::bindings::http::types::Method::Head => {
                req.method(http::Method::HEAD)
            }
            wasmtime_wasi_http::bindings::http::types::Method::Post => {
                req.method(http::Method::POST)
            }
            wasmtime_wasi_http::bindings::http::types::Method::Put => req.method(http::Method::PUT),
            wasmtime_wasi_http::bindings::http::types::Method::Delete => {
                req.method(http::Method::DELETE)
            }
            wasmtime_wasi_http::bindings::http::types::Method::Connect => {
                req.method(http::Method::CONNECT)
            }
            wasmtime_wasi_http::bindings::http::types::Method::Options => {
                req.method(http::Method::OPTIONS)
            }
            wasmtime_wasi_http::bindings::http::types::Method::Trace => {
                req.method(http::Method::TRACE)
            }
            wasmtime_wasi_http::bindings::http::types::Method::Patch => {
                req.method(http::Method::PATCH)
            }
            wasmtime_wasi_http::bindings::http::types::Method::Other(other) => {
                req.method(other.as_str())
            }
        };
        let req = req
            .uri(uri)
            .body(())
            .context("failed to build HTTP request")?;
        let (parts, ()) = req.into_parts();

        let req = HostIncomingRequest::new(
            self,
            parts,
            scheme.unwrap_or(wasmtime_wasi_http::bindings::http::types::Scheme::Http),
            body.map(|body| {
                wasmtime_wasi_http::body::HostIncomingBody::new(body, Duration::from_secs(1))
            }),
        )
        .context("failed to construct `incoming-request`")?;
        WasiHttpView::table(self)
            .push(req)
            .context("failed to push `incoming-request` into resource table")
    }
}

pub struct Config<'a> {
    pub engine: Engine,
    pub wasm: &'a [u8],
}

impl Default for Config<'_> {
    fn default() -> Self {
        Self {
            engine: Engine::default(),
            wasm: PASSTHROUGH,
        }
    }
}

pub struct Func<'a> {
    func: wasmtime::component::Func,
    store: &'a mut Store<Ctx>,
}

impl Func<'_> {
    pub fn params(&self) -> Box<[Type]> {
        self.func.params(&self.store)
    }

    pub fn results(&self) -> Box<[Type]> {
        self.func.results(&self.store)
    }

    pub fn call(&mut self, params: &[Val], results: &mut [Val]) -> anyhow::Result<()> {
        self.func
            .call(&mut self.store, params, results)
            .context("failed to call function")?;
        self.func
            .post_return(&mut self.store)
            .context("failed to invoke `post-return`")
    }

    pub fn store(&mut self) -> &mut Store<impl WasiView + WasiHttpView> {
        &mut self.store
    }
}

pub struct Instance {
    instance: wasmtime::component::Instance,
    store: Store<Ctx>,
}

impl Instance {
    pub fn func(&mut self, instance: &str, name: &str) -> anyhow::Result<Func> {
        let idx = self
            .instance
            .get_export(&mut self.store, None, instance)
            .with_context(|| format!("export `{instance}` not found"))?;
        let idx = self
            .instance
            .get_export(&mut self.store, Some(&idx), name)
            .with_context(|| format!("export `{name}` not found"))?;
        let func = self
            .instance
            .get_func(&mut self.store, idx)
            .with_context(|| format!("function export `{name}` not found"))?;
        Ok(Func {
            func,
            store: &mut self.store,
        })
    }

    pub fn call(
        &mut self,
        instance: &str,
        name: &str,
        params: &[Val],
        results: &mut [Val],
    ) -> anyhow::Result<()> {
        let mut func = self
            .func(instance, name)
            .context("failed to lookup function")?;
        func.call(params, results)
            .context("failed to call function")
    }

    pub fn call_http_response_outparam_set(
        &mut self,
        out: Resource<wasmtime_wasi_http::types::HostResponseOutparam>,
        res: Result<
            Resource<wasmtime_wasi_http::types::HostOutgoingResponse>,
            wasmtime_wasi_http::bindings::http::types::ErrorCode,
        >,
    ) -> anyhow::Result<()> {
        let func = self
            .func("wasi:http/types@0.2.0", "[static]response-outparam.set")
            .context("failed to lookup function")?;
        let func = unsafe { TypedFunc::new_unchecked(func.func) };
        func.call(&mut self.store, (out.rep(), res))
            .context("failed to call function")
    }

    pub fn store(&mut self) -> &mut Store<impl WasiView + WasiHttpView> {
        &mut self.store
    }
}

pub fn instantiate(Config { engine, wasm }: Config) -> anyhow::Result<Instance> {
    let wasm = if wasmparser::Parser::is_core_wasm(&wasm) {
        let wasm = wit_component::ComponentEncoder::default()
            .module(&wasm)
            .context("failed to set core component module")?
            .adapter(
                WASI_SNAPSHOT_PREVIEW1_ADAPTER_NAME,
                WASI_SNAPSHOT_PREVIEW1_REACTOR_ADAPTER,
            )
            .context("failed to add WASI preview1 adapter")?
            .encode()
            .context("failed to encode a component from module")?;
        let wasm = wasm.as_slice();
        return instantiate(Config { engine, wasm });
    } else {
        wasm
    };
    let wasm = Component::new(&engine, wasm).context("failed to compile component")?;

    let mut linker = Linker::<Ctx>::new(&engine);
    wasmtime_wasi::add_to_linker_sync(&mut linker).context("failed to link WASI")?;
    wasmtime_wasi_http::add_only_http_to_linker_sync(&mut linker)
        .context("failed to link `wasi:http`")?;
    bindings::west::test::http_test::add_to_linker(&mut linker, |cx| cx)?;

    let wasi = WasiCtxBuilder::new()
        .inherit_env()
        .inherit_stdout()
        .inherit_stderr()
        .inherit_network()
        .build();
    let http = WasiHttpCtx::new();
    let table = ResourceTable::new();
    let mut store = Store::new(&engine, Ctx { wasi, http, table });
    let instance = linker
        .instantiate(&mut store, &wasm)
        .context("failed to instantiate component")?;
    Ok(Instance { instance, store })
}
