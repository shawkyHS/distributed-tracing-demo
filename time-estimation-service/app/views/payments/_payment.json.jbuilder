json.extract! payment, :id, :amount, :created_at, :updated_at
json.url payment_url(payment, format: :json)
