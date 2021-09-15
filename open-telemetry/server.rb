#!/usr/bin/env ruby
# frozen_string_literal: true

# Copyright The OpenTelemetry Authors
#
# SPDX-License-Identifier: Apache-2.0

#require 'rubygems'
#require 'bundler/setup'
require 'sinatra'
# Require otel-ruby
require 'opentelemetry/sdk'
require 'opentelemetry/exporter/jaeger'

# Export traces to console by default (exporter)
ENV['OTEL_TRACES_EXPORTER'] ||= 'jaeger' # can be console

# configure SDK with defaults, # configure for console 
OpenTelemetry::SDK.configure

OpenTelemetry::SDK.configure do |c|
  # c.add_span_processor(
  #   OpenTelemetry::SDK::Trace::Export::BatchSpanProcessor.new(
  #     exporter: OpenTelemetry::Exporter::Jaeger::AgentExporter.new(host: '127.0.0.1', port: 6831)
  #     # Alternatively, for the collector exporter:
  #     # exporter: OpenTelemetry::Exporter::Jaeger::CollectorExporter.new(endpoint: 'http://192.168.0.1:14268/api/traces')
  #   )
  # )
  c.service_name = 'open-telemetry-example-server'
  c.service_version = '0.1.0'
end

# Rack middleware to extract span context, create child span, and add
# attributes/events to the span
class OpenTelemetryMiddleware
  def initialize(app)
    @app = app
    @tracer = OpenTelemetry.tracer_provider.tracer('rails', '1.0')
  end

  def call(env)
    # Extract context from request headers
    context = OpenTelemetry.propagation.extract(
      env,
      getter: OpenTelemetry::Context::Propagation.rack_env_getter
    )

    status, headers, response_body = 200, {}, ''

    # Span name SHOULD be set to route:
    span_name = env['PATH_INFO']

    # For attribute naming, see
    # https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/data-semantic-conventions.md#http-server

    # Activate the extracted context
    OpenTelemetry::Context.with_current(context) do
      # Span kind MUST be `:server` for a HTTP server span
      @tracer.in_span(
        span_name,
        attributes: {
          'component' => 'http',
          'http.method' => env['REQUEST_METHOD'],
          'http.route' => env['PATH_INFO'],
          'http.url' => env['REQUEST_URI'],
        },
        kind: :server
      ) do |span|
        # Run application stack
        status, headers, response_body = @app.call(env)

        span.set_attribute('http.status_code', status)
        span.add_event('test_event', attributes: { user_id: '123'})
      end
    end

    [status, headers, response_body]
  end
end

class App < Sinatra::Base
  set :bind, '0.0.0.0'
  use OpenTelemetryMiddleware

  get '/hello' do
    'Hello World!'
  end

  run! if app_file == $0
end
