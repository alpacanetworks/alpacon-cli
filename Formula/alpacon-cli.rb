# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class AlpaconCli < Formula
  desc "CLI for Alpacon"
  homepage "https://github.com/alpacanetworks/alpacon-cli"
  version "0.2.1"

  on_macos do
    on_intel do
      url "https://github.com/alpacanetworks/alpacon-cli/releases/download/0.2.1/alpacon-0.2.1-darwin-amd64.tar.gz"
      sha256 "fdbd6dc159d97b1560866c23f0e66b98a7da785e319ed31f3cff32221a42b735"

      def install
        bin.install "alpacon"
      end
    end
    on_arm do
      url "https://github.com/alpacanetworks/alpacon-cli/releases/download/0.2.1/alpacon-0.2.1-darwin-arm64.tar.gz"
      sha256 "62b7d09b7cdf276c5a0e36b7141c94bc8b6691de14a0dbd013e559336c8edeb0"

      def install
        bin.install "alpacon"
      end
    end
  end

  on_linux do
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/alpacanetworks/alpacon-cli/releases/download/0.2.1/alpacon-0.2.1-linux-amd64.tar.gz"
        sha256 "cedd3509683fcff51436cada37df36c215feb1dce76ad10a27432b30704bf419"

        def install
          bin.install "alpacon"
        end
      end
    end
    on_arm do
      if !Hardware::CPU.is_64_bit?
        url "https://github.com/alpacanetworks/alpacon-cli/releases/download/0.2.1/alpacon-0.2.1-linux-arm.tar.gz"
        sha256 "99ca7ed9124b5dfb9c16aa2228ca34073634b24bc226870f5a1b2ab582e28944"

        def install
          bin.install "alpacon"
        end
      end
    end
    on_arm do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/alpacanetworks/alpacon-cli/releases/download/0.2.1/alpacon-0.2.1-linux-arm64.tar.gz"
        sha256 "b4c9190994f18c3922c36b76c71678d2275ea12c302e1ab8f025605fc88c57e9"

        def install
          bin.install "alpacon"
        end
      end
    end
  end
end
