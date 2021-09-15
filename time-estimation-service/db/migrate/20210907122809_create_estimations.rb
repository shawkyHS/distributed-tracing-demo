class CreateEstimations < ActiveRecord::Migration[6.0]
  def change
    create_table :estimations do |t|
      t.integer :minutes

      t.timestamps
    end
  end
end
