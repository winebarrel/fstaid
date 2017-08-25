require 'formula'

class Fstaid < Formula
  VERSION = '0.1.6'

  homepage 'https://github.com/winebarrel/fstaid'
  url "https://github.com/winebarrel/fstaid/releases/download/v#{VERSION}/fstaid-v#{VERSION}-darwin-amd64.gz"
  sha256 '9a816ce467c1d999a5795cc317fb096bf5c2b10875c7c400c2013d5fc637e7f8'
  version VERSION
  head 'https://github.com/winebarrel/fstaid.git', :branch => 'master'

  def install
    system "mv fstaid-v#{VERSION}-darwin-amd64 fstaid"
    bin.install 'fstaid'
  end
end
