Rails.application.routes.draw do
  resources :estimations
  resources :groceries
  resources :payments
  resources :availabilities
  resources :orders
  resources :users
  # For details on the DSL available within this file, see https://guides.rubyonrails.org/routing.html
end
