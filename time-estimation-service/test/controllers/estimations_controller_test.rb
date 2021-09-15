require 'test_helper'

class EstimationsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @estimation = estimations(:one)
  end

  test "should get index" do
    get estimations_url
    assert_response :success
  end

  test "should get new" do
    get new_estimation_url
    assert_response :success
  end

  test "should create estimation" do
    assert_difference('Estimation.count') do
      post estimations_url, params: { estimation: { minutes: @estimation.minutes } }
    end

    assert_redirected_to estimation_url(Estimation.last)
  end

  test "should show estimation" do
    get estimation_url(@estimation)
    assert_response :success
  end

  test "should get edit" do
    get edit_estimation_url(@estimation)
    assert_response :success
  end

  test "should update estimation" do
    patch estimation_url(@estimation), params: { estimation: { minutes: @estimation.minutes } }
    assert_redirected_to estimation_url(@estimation)
  end

  test "should destroy estimation" do
    assert_difference('Estimation.count', -1) do
      delete estimation_url(@estimation)
    end

    assert_redirected_to estimations_url
  end
end
