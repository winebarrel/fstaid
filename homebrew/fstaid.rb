require 'formula'

class Fstaid < Formula
  VERSION = '0.1.4'

  homepage 'https://github.com/winebarrel/fstaid'
  url "https://github.com/winebarrel/fstaid/releases/download/v#{VERSION}/fstaid-v#{VERSION}-darwin-amd64.gz"
  sha256 'b4e91d3b04240d21a2a9149e6206dbae1bc5ac6729ad668378da4fea47b6955b'
  version VERSION
  head 'https://github.com/winebarrel/fstaid.git', :branch => 'master'

  def install
    system "mv fstaid-v#{VERSION}-darwin-amd64 fstaid"
    bin.install 'fstaid'
  end
end
