%define  debug_package %{nil}

Name:     fstaid
Version:  0.1.4
Release:  1%{?dist}
Summary:  fstaid is a daemon that monitors the health condition of the server and executes the script if there is any problem.

Group:    System Environment/Daemons
License:  MIT
URL:    https://github.com/winebarrel/fstaid
Source0:  %{name}.tar.gz
# https://github.com/winebarrel/fstaid/releases/download/v%{version}/fstaid_%{version}.tar.gz

%description
fstaid is a daemon that monitors the health condition of the server and executes the script if there is any problem.

%prep
%setup -q -n src

%build
make
make test

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/usr/sbin
install -m 755 fstaid %{buildroot}/usr/sbin/

%files
%defattr(755,root,root,-)
/usr/sbin/fstaid
