require "google/cloud/pubsub"

class GroceriesController < ApplicationController
  before_action :set_grocery, only: %i[ show edit update destroy ]

  # GET /groceries or /groceries.json
  def index
    sleep 1
    fetch_time_estimatiom
    transmit_order
    update_order_status
    @groceries = Grocery.all
  end

  # GET /groceries/1 or /groceries/1.json
  def show
  end

  # GET /groceries/new
  def new
    @grocery = Grocery.new
  end

  # GET /groceries/1/edit
  def edit
  end

  # POST /groceries or /groceries.json
  def create
    @grocery = Grocery.new(grocery_params)

    respond_to do |format|
      if @grocery.save
        format.html { redirect_to @grocery, notice: "Grocery was successfully created." }
        format.json { render :show, status: :created, location: @grocery }
      else
        format.html { render :new, status: :unprocessable_entity }
        format.json { render json: @grocery.errors, status: :unprocessable_entity }
      end
    end
  end

  # PATCH/PUT /groceries/1 or /groceries/1.json
  def update
    respond_to do |format|
      if @grocery.update(grocery_params)
        format.html { redirect_to @grocery, notice: "Grocery was successfully updated." }
        format.json { render :show, status: :ok, location: @grocery }
      else
        format.html { render :edit, status: :unprocessable_entity }
        format.json { render json: @grocery.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /groceries/1 or /groceries/1.json
  def destroy
    @grocery.destroy
    respond_to do |format|
      format.html { redirect_to groceries_url, notice: "Grocery was successfully destroyed." }
      format.json { head :no_content }
    end
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_grocery
      @grocery = Grocery.find(params[:id])
    end

    # Only allow a list of trusted parameters through.
    def grocery_params
      params.require(:grocery).permit(:amount)
    end

    def fetch_time_estimatiom
      Faraday.get 'http://localhost:3004/estimations'
    end

    def transmit_order
        OpenTelemetry.tracer_provider.tracer.in_span('transmit_order') do |s|
          sleep 1
          s.add_event('order_transmitted', { attributes: { 'remote_id' => '1010'}})
          s.add_attributes({'name' => 'Ahmed', 'order_id' => 5})
          #s.status = OpenTelemetry::Trace::Status::OK
        end
    end

    def update_order_status
      pubsub = Google::Cloud::PubSub.new(
        project_id: "pubsub-quick-sta-1625046790078", 
        credentials: "/Users/ahmedshawky/learn/pubsub/ruby/pubsub-quick-sta-1625046790078-8e656f44ad6d.json"
      )
      topic_id = "first-topic"
      topic = pubsub.topic topic_id

      OpenTelemetry.tracer_provider.tracer.in_span('update_order_status', attributes: {}, kind: :client) do |s|

        data = {}
        OpenTelemetry.propagation.inject(data) # inject the context
        status = OrderStatus.new(id: 10, key: "received_by_vendor", trace_context: data.to_json)

        topic.publish_async OrderStatus.encode(status) do |result|
          raise "Failed to publish the message." unless result.succeeded?
          puts "Message published asynchronously."
        end

      end
    end

end
