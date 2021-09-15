require "application_system_test_case"

class GroceriesTest < ApplicationSystemTestCase
  setup do
    @grocery = groceries(:one)
  end

  test "visiting the index" do
    visit groceries_url
    assert_selector "h1", text: "Groceries"
  end

  test "creating a Grocery" do
    visit groceries_url
    click_on "New Grocery"

    fill_in "Amount", with: @grocery.amount
    click_on "Create Grocery"

    assert_text "Grocery was successfully created"
    click_on "Back"
  end

  test "updating a Grocery" do
    visit groceries_url
    click_on "Edit", match: :first

    fill_in "Amount", with: @grocery.amount
    click_on "Update Grocery"

    assert_text "Grocery was successfully updated"
    click_on "Back"
  end

  test "destroying a Grocery" do
    visit groceries_url
    page.accept_confirm do
      click_on "Destroy", match: :first
    end

    assert_text "Grocery was successfully destroyed"
  end
end
