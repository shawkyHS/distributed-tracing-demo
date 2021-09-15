class FinanceWorker
  include Sidekiq::Worker

  def perform
    sleep 1
    #raise NotImplementedError
  end
end
