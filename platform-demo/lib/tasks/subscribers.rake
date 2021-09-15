require "google/cloud/pubsub"

namespace :subscribers do
  desc "TODO"
  task order_status: :environment do

    tracer = OpenTelemetry.tracer_provider.tracer('platform_order_status_pubsub_subscription', '1.0')

    pubsub = Google::Cloud::PubSub.new(
      project_id: "pubsub-quick-sta-1625046790078",
      credentials: "/Users/ahmedshawky/learn/pubsub/ruby/pubsub-quick-sta-1625046790078-8e656f44ad6d.json"
    )
    sub = pubsub.subscription "platform_order_status_subscription"
    subscriber = sub.listen do |received_message|
      order_status = OrderStatus.decode(received_message.message.data)

      context = OpenTelemetry.propagation.extract(
        JSON.parse(order_status.trace_context),
        getter: OpenTelemetry::Context::Propagation.text_map_getter
      )

      OpenTelemetry::Context.with_current(context) do
        tracer.in_span('order_status_subscriber') do |span|
          span.add_event('order_status_received', { attributes: { 'status_key' => order_status.key}})
          puts "Data: #{received_message.message.data}, published at #{received_message.message.published_at}"
          received_message.acknowledge!
          sleep 0.5
        end
      end
    end

    # Handle exceptions from listener
    subscriber.on_error do |exception|
      puts "Exception: #{exception.class} #{exception.message}"
    end

    at_exit do
      subscriber.stop!(10)
    end

    subscriber.start

    sleep

  end

end
