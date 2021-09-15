#!/usr/bin/env ruby
# frozen_string_literal: true

# Copyright The OpenTelemetry Authors
#
# SPDX-License-Identifier: Apache-2.0
#require 'bundler/setup'
require 'faraday'
# Require otel-ruby
require 'opentelemetry/sdk'
require 'opentelemetry/exporter/jaeger'
require 'byebug'

# Export traces to console by default
ENV['OTEL_TRACES_EXPORTER'] ||= 'jaeger' # can be console

# Allow setting the host from the ENV
host = ENV.fetch('HTTP_EXAMPLE_HOST', '0.0.0.0')

# configure SDK with defaults
OpenTelemetry::SDK.configure

OpenTelemetry::SDK.configure do |c|
  # c.add_span_processor(
  #   OpenTelemetry::SDK::Trace::Export::BatchSpanProcessor.new(
  #     exporter: OpenTelemetry::Exporter::Jaeger::AgentExporter.new(host: '127.0.0.1', port: 6831)
  #     # Alternatively, for the collector exporter:
  #     # exporter: OpenTelemetry::Exporter::Jaeger::CollectorExporter.new(endpoint: 'http://192.168.0.1:14268/api/traces')
  #   )
  # )
  c.service_name = 'open-telemetry-example-client'
  c.service_version = '0.1.0'
end

# Configure tracer
tracer = OpenTelemetry.tracer_provider.tracer('faraday', '1.0')

connection = Faraday.new("http://#{host}:4567")
url = '/hello'

# For attribute naming, see:
# https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/data-semantic-conventions.md#http-client

# Span name should be set to URI path value:
tracer.in_span(
  url,
  attributes: {
    'component' => 'http',
    'http.method' => 'GET',
  },
  kind: :client
) do |span|
  response = connection.get(url) do |request|
    # Inject context into request headers
    OpenTelemetry.propagation.inject(request.headers)
  end

  span.set_attribute('http.url', response.env.url.to_s)
  span.set_attribute('http.status_code', response.status)
end
