# Load the Rails application.
require_relative 'application'

# Initialize the Rails application.
Rails.application.initialize!

logger = Logger.new(STDOUT)
logger.progname = 'multivac'
original_formatter = Logger::Formatter.new
logger.formatter  = proc do |severity, datetime, progname, msg|
  current_span = OpenTelemetry::Trace.current_span(OpenTelemetry::Context.current).context
  
  dd_trace_id = current_span.trace_id.unpack1('H*')[16, 16].to_i(16).to_s
  dd_span_id = current_span.span_id.unpack1('H*').to_i(16).to_s
  
  if current_span
    "#{{datetime: datetime, progname: progname, severity: severity, msg: msg, 'dd.trace_id': dd_trace_id, 'dd.span_id': dd_span_id}.to_json}\n"
  else
    "#{{datetime: datetime, progname: progname, severity: severity, msg: msg}.to_json}\n"
  end
end

logger.info("Example log line with trace correlation info")
