class EstimationsController < ApplicationController
  before_action :set_estimation, only: %i[ show edit update destroy ]

  # GET /estimations or /estimations.json
  def index
    sleep 1
    span = OpenTelemetry::Trace.current_span
    span.add_event('payload_data', { attributes: { 'pickup_time' => '2 sec'}})
    span.add_attributes({'name' => 'Ahmed', 'order_id' => 5})
    # span.status = OpenTelemetry::Trace::Status::ERROR
    @estimations = Estimation.all
  end

  # GET /estimations/1 or /estimations/1.json
  def show
  end

  # GET /estimations/new
  def new
    @estimation = Estimation.new
  end

  # GET /estimations/1/edit
  def edit
  end

  # POST /estimations or /estimations.json
  def create
    @estimation = Estimation.new(estimation_params)

    respond_to do |format|
      if @estimation.save
        format.html { redirect_to @estimation, notice: "Estimation was successfully created." }
        format.json { render :show, status: :created, location: @estimation }
      else
        format.html { render :new, status: :unprocessable_entity }
        format.json { render json: @estimation.errors, status: :unprocessable_entity }
      end
    end
  end

  # PATCH/PUT /estimations/1 or /estimations/1.json
  def update
    respond_to do |format|
      if @estimation.update(estimation_params)
        format.html { redirect_to @estimation, notice: "Estimation was successfully updated." }
        format.json { render :show, status: :ok, location: @estimation }
      else
        format.html { render :edit, status: :unprocessable_entity }
        format.json { render json: @estimation.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /estimations/1 or /estimations/1.json
  def destroy
    @estimation.destroy
    respond_to do |format|
      format.html { redirect_to estimations_url, notice: "Estimation was successfully destroyed." }
      format.json { head :no_content }
    end
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_estimation
      @estimation = Estimation.find(params[:id])
    end

    # Only allow a list of trusted parameters through.
    def estimation_params
      params.require(:estimation).permit(:minutes)
    end
end
