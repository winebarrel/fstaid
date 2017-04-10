require 'formula'

class Fstaid < Formula
  VERSION = '0.1.1'

  homepage 'https://github.com/winebarrel/fstaid'
  url "https://github.com/winebarrel/fstaid/releases/download/v#{VERSION}/fstaid-v#{VERSION}-darwin-amd64.gz"
  sha256 '16a93005c30dc54640512de1635959fa94507e1b45e0af3ef90fa847a9618554'
  version VERSION
  head 'https://github.com/winebarrel/fstaid.git', :branch => 'master'

  def install
    system "mv fstaid-v#{VERSION}-darwin-amd64 fstaid"
    bin.install 'fstaid'
  end
end
