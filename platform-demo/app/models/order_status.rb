# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: order_status.proto

require 'google/protobuf'

Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("order_status.proto", :syntax => :proto3) do
    add_message "OrderStatus" do
      optional :id, :int32, 1
      optional :key, :string, 2
      optional :trace_context, :string, 3
    end
    add_message "OrderStatuses" do
      repeated :order_statuses, :message, 1, "OrderStatus"
    end
  end
end

OrderStatus = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("OrderStatus").msgclass
#OrderStatuses = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("OrderStatuses").msgclass
