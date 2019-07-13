Name: pet-adoption
Version:1.0.0
Release:1
Summary:pet-adoption

#Group:		
License:Commercial
#URL:		
#Source0:	
#BuildRoot:	%(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

#BuildRequires:	

#%global __os_install_post %{nil}
%description
pet-adoption

%prep

%build

%install
cp $RPM_BUILD_DIR/%{name}/* $RPM_BUILD_ROOT/ -raf

%clean
rm -rf %{buildroot}

%pre


%post
systemctl daemon-reload

%preun


%postun

%files
%defattr(-,root,root,-)
/usr/local/bin/pet-adoption
/etc/systemd/system/pet-adoption.service
%config(noreplace) /etc/pet-adoption/config.yml

%changelog

