# ENV['OTEL_TRACES_EXPORTER'] ||= 'jaeger'
require 'opentelemetry/exporter/otlp'
# ENV['OTEL_TRACES_EXPORTER'] = 'otlp'
ENV['OTEL_EXPORTER_OTLP_ENDPOINT'] = 'http://localhost:4318'
ENV['OTEL_RUBY_EXPORTER_OTLP_SSL_VERIFY_NONE'] = 'OpenSSL::SSL:VERIFY_NONE'

# ENV['OTEL_RUBY_EXPORTER_OTLP_SSL_VERIFY_PEER']= 'OpenSSL::SSL:VERIFY_NONE'
OpenTelemetry::SDK.configure do |c|
  c.service_name = 'order-transmission-microservice'
  c.service_version = '0.1.0'
  c.use_all({
    'OpenTelemetry::Instrumentation::Sidekiq' => { span_naming: :job_class, propagation_style: :child },
    'OpenTelemetry::Instrumentation::Rails' =>   { enable_recognize_route: true }
    })
end

# #puts OpenTelemetry::Instrumentation::Sidekiq::Instrumentation.instance.name
# OpenTelemetry.tracer_provider.tracer
# OpenTelemetry::Instrumentation::Sidekiq

