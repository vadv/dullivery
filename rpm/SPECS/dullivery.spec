%define version unknown
%define debug_package %{nil}
Name:           dullivery
Version:        %{version}
Release:        1%{?dist}
Summary:        dullivery
License:        BSD
URL:            http://git.itv.restr.im/infra/dullivery
Source1:        dullivery.init
Source2:        dullivery-web.init
Source3:        dullivery-logrotate.in
Source:         dullivery-%{version}.tar.gz
BuildRequires:  make
BuildRequires:  git

%define restream_dir /opt/restream/
%define restream_bin_dir %{restream_dir}/dullivery/bin
%define dullivery_static_dir %{restream_dir}/dullivery/share
%define dullivery_data_dir /var/lib/dullivery

%description
This package provides dullivery system.

%prep
%setup

%pre

getent group dullivery > /dev/null || groupadd -r dullivery
getent passwd dullivery > /dev/null || \
    useradd -r -g dullivery -d /var/run/dullivery -s /sbin/nologin \
    -c "dullivery user" dullivery

mkdir -p /var/log/dullivery
chown dullivery:dullivery /var/log/dullivery

mkdir -p /var/lib/dullivery
chown dullivery:dullivery /var/lib/dullivery

%build
make

%install
%{__mkdir} -p %{buildroot}%{restream_bin_dir}
%{__install} -m 0755 -p bin/dullivery %{buildroot}%{restream_bin_dir}
%{__install} -m 0755 -p bin/dullivery-web %{buildroot}%{restream_bin_dir}
%{__mkdir} -p %{buildroot}/%{_sysconfdir}/init.d
%{__mkdir} -p %{buildroot}/%{_sysconfdir}/logrotate.d
%{__install} -m 0755 -p %{SOURCE1} %{buildroot}/%{_sysconfdir}/init.d/dullivery
%{__install} -m 0755 -p %{SOURCE2} %{buildroot}/%{_sysconfdir}/init.d/dullivery-web
%{__install} -m 0644 -p %{SOURCE3} %{buildroot}/%{_sysconfdir}/logrotate.d/dullivery
%{__mkdir} -p %{buildroot}%{restream_dir}/dullivery
cp -rva src/web/public %{buildroot}%{dullivery_static_dir}
%{__mkdir} -p %{buildroot}%{dullivery_data_dir}

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%{restream_bin_dir}/dullivery
%{restream_bin_dir}/dullivery-web
%{dullivery_static_dir}/*
%doc README.md
%{_sysconfdir}/init.d/dullivery
%{_sysconfdir}/init.d/dullivery-web
%{_sysconfdir}/logrotate.d/dullivery
