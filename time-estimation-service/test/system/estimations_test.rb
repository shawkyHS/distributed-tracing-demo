require "application_system_test_case"

class EstimationsTest < ApplicationSystemTestCase
  setup do
    @estimation = estimations(:one)
  end

  test "visiting the index" do
    visit estimations_url
    assert_selector "h1", text: "Estimations"
  end

  test "creating a Estimation" do
    visit estimations_url
    click_on "New Estimation"

    fill_in "Minutes", with: @estimation.minutes
    click_on "Create Estimation"

    assert_text "Estimation was successfully created"
    click_on "Back"
  end

  test "updating a Estimation" do
    visit estimations_url
    click_on "Edit", match: :first

    fill_in "Minutes", with: @estimation.minutes
    click_on "Update Estimation"

    assert_text "Estimation was successfully updated"
    click_on "Back"
  end

  test "destroying a Estimation" do
    visit estimations_url
    page.accept_confirm do
      click_on "Destroy", match: :first
    end

    assert_text "Estimation was successfully destroyed"
  end
end
