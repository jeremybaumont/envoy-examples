use log::{debug, error};
use proxy_wasm::traits::*;
use proxy_wasm::types::*;

#[no_mangle]
pub fn _start() {
    proxy_wasm::set_log_level(LogLevel::Debug);
    proxy_wasm::set_root_context(|_| -> Box<dyn RootContext> { Box::new(ExtAuthzMismatchRoot) });
}

struct ExtAuthzMismatchRoot;

impl Context for ExtAuthzMismatchRoot {}

impl RootContext for ExtAuthzMismatchRoot {
    fn get_type(&self) -> Option<ContextType> {
        Some(ContextType::HttpContext)
    }

    fn create_http_context(&self, context_id: u32) -> Option<Box<dyn HttpContext>> {
        Some(Box::new(ExtAuthzMismatch::new(context_id)))
    }
}


#[derive(Debug)]
struct ExtAuthzMismatch {
    context_id: u32,
    extauthz_status: String,
    auth_mismatch_metric_id: u32,
}

impl Context for ExtAuthzMismatch {}

impl ExtAuthzMismatch {
    fn new(context_id: u32) -> Self {
        return Self {
            context_id,
            extauthz_status: String::from(""),
            auth_mismatch_metric_id: proxy_wasm::hostcalls::define_metric(MetricType::Counter, "authMismatch").unwrap(),
        }
    }
}

impl HttpContext for ExtAuthzMismatch {

    fn on_http_request_headers(&mut self, _num_headers: usize) -> Action {
        self.extauthz_status = self.get_http_request_header("x-extauthz-status-code").unwrap_or(String::from(""));
        Action::Continue
    }

    fn on_http_response_headers(&mut self, _num_headers: usize) -> Action {
        let resp_status = self.get_http_response_header(":status").unwrap_or(String::from(""));

        debug!("context_id : {:?}, extauthz status code : {:?}, response status code : {:?}", self.context_id, self.extauthz_status, resp_status);
        if resp_status != self.extauthz_status && self.extauthz_status != "" {
            proxy_wasm::hostcalls::increment_metric(self.auth_mismatch_metric_id, 1);
        }
        Action::Continue
    }

}
