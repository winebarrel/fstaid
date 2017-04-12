require 'formula'

class Fstaid < Formula
  VERSION = '0.1.2'

  homepage 'https://github.com/winebarrel/fstaid'
  url "https://github.com/winebarrel/fstaid/releases/download/v#{VERSION}/fstaid-v#{VERSION}-darwin-amd64.gz"
  sha256 '27089a6b6dcdd0fbfe3c25e05909fb8b892afa4b6d25bdda144b0aebf4c3c6ab'
  version VERSION
  head 'https://github.com/winebarrel/fstaid.git', :branch => 'master'

  def install
    system "mv fstaid-v#{VERSION}-darwin-amd64 fstaid"
    bin.install 'fstaid'
  end
end
