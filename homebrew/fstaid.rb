require 'formula'

class Fstaid < Formula
  VERSION = '0.1.7'

  homepage 'https://github.com/winebarrel/fstaid'
  url "https://github.com/winebarrel/fstaid/releases/download/v#{VERSION}/fstaid-v#{VERSION}-darwin-amd64.gz"
  sha256 'd9a16aca51cd921325aa4eb42c3e87b99ab960f7387a0877eca0f8c3e3a0a8a9'
  version VERSION
  head 'https://github.com/winebarrel/fstaid.git', :branch => 'master'

  def install
    system "mv fstaid-v#{VERSION}-darwin-amd64 fstaid"
    bin.install 'fstaid'
  end
end
