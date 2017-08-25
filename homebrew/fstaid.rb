require 'formula'

class Fstaid < Formula
  VERSION = '0.1.5'

  homepage 'https://github.com/winebarrel/fstaid'
  url "https://github.com/winebarrel/fstaid/releases/download/v#{VERSION}/fstaid-v#{VERSION}-darwin-amd64.gz"
  sha256 '61710717a0e5f8092bcbfd67b6fbc8a418eae6d6ccd266d6391d8b16afe9c7a1'
  version VERSION
  head 'https://github.com/winebarrel/fstaid.git', :branch => 'master'

  def install
    system "mv fstaid-v#{VERSION}-darwin-amd64 fstaid"
    bin.install 'fstaid'
  end
end
