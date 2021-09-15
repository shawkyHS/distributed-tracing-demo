class CreateGroceries < ActiveRecord::Migration[6.0]
  def change
    create_table :groceries do |t|
      t.string :amount

      t.timestamps
    end
    add_index :groceries, :amount
  end
end
