class OpenTelemetryMiddleware
  def initialize(app)
    @app = app
    @tracer = OpenTelemetry.tracer_provider.tracer('time-estimation-service', '1.0')
  end

  def call(env)
    # Extract context from request headers
    context = OpenTelemetry.propagation.extract(
      env,
      getter: OpenTelemetry::Context::Propagation.rack_env_getter
    )

    status, headers, response_body = 200, {}, ''

    # For attribute naming, see
    # https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/data-semantic-conventions.md#http-server

    # Activate the extracted context
    OpenTelemetry::Context.with_current(context) do
      # Span kind MUST be `:server` for a HTTP server span
      @tracer.in_span(
        env['PATH_INFO'], # span name set to the route
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
        #span.add_event('test_event', attributes: { user_id: '123'})
      end
    end

    [status, headers, response_body]
  end
end
