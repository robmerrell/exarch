defmodule TestApp.Accounts.Users do
  @moduledoc """
  A fake module for unit test
  """

  alias TestApp.Repo
  alias TestApp.Accounts.{User, Admin}

  @doc """
  Get a user by id
  """
  def get_user!(id), do: Repo.get!(User, id)

  def get_by_username(username) do
    Repo.get_by(User, username: username)
  end

  def update_user(user, attrs) do
    user
    |> User.changeset(attrs)
    |> Repo.update()
  end

  @doc """
  Returns a hello message for the user as a string
  """
  defp hello_message(user) do
    color =
      case user.favorite_fruit do
        :blueberry -> "blue"
        :strawberry -> "red"
        :lime -> "green"
      end

    message = "Hello #{user.first}, your favorite color is: #{color}"
  end

  def function_without_args do
    [
      "string one",
      "string two",
      "string three"
    ]
  end
end
