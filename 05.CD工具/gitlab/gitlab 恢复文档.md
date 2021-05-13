# 环境
1. 原gitlab 版本11.6.3-ee (76c33c6)，操作系统： centos7.4
 
 
## 新环境初始化
1. 运行脚本下面的 shell脚本 `sh script.rpm.sh` 脚本在最后
2. 安装gitlab ` yum -y install gitlab-ee-11.6.3-ee.0.el7.x86_64 ` 版本不知道可以在web界面查看，或者查看yum.log安装历史
3. 备份服务器上的配置文件`cp  /etc/gitlab/gitlab.rb /etc/gitlab/gitlab.rb.bak`
4. 原服务器配置文件拷贝过去进行替换  `\cp gitlab.rb /etc/gitlab/gitlab.rb`
5. gitlab-ctl reconfig && gitlab-ctl restart 
已上就是配置好之前的gitlab.rb配置文件，和重新配置启动gitlab

### 源服务器备份数据
1. 运行备份命令 `/usr/bin/gitlab-rake gitlab:backup:create`
2. 自动生成备份文件，手动上传到新服务器  `1620630804_2021_05_10_11.6.3-ee_gitlab_backup.tar` 路径按照配置文件提示地方
```
[root@gitlab ~]# cat /etc/gitlab/gitlab.rb| grep backup_path
gitlab_rails['manage_backup_path'] = true
gitlab_rails['backup_path'] = "/data/gitlab_backup" 
```

### 新服务器上应用归档文件
1. `gitlab-rake gitlab:backup:restore BACKUP=1620630804_2021_05_10_11.6.3-ee`
2. 等待一会需要输入yes 就 ok了
3. 访问新服务器进行验证


#### script.rpm.sh 脚本
```
#!/bin/bash

unknown_os ()
{
  echo "Unfortunately, your operating system distribution and version are not supported by this script."
  echo
  echo "You can override the OS detection by setting os= and dist= prior to running this script."
  echo "You can find a list of supported OSes and distributions on our website: https://packages.gitlab.com/docs#os_distro_version"
  echo
  echo "For example, to force CentOS 6: os=el dist=6 ./script.sh"
  echo
  echo "Please email support@packagecloud.io and let us know if you run into any issues."
  exit 1
}

curl_check ()
{
  echo "Checking for curl..."
  if command -v curl > /dev/null; then
    echo "Detected curl..."
  else
    echo "Installing curl..."
    yum install -d0 -e0 -y curl
  fi
}


detect_os ()
{
  if [[ ( -z "${os}" ) && ( -z "${dist}" ) ]]; then
    if [ -e /etc/os-release ]; then
      . /etc/os-release
      os=${ID}
      if [ "${os}" = "poky" ]; then
        dist=`echo ${VERSION_ID}`
      elif [ "${os}" = "sles" ]; then
        dist=`echo ${VERSION_ID}`
      elif [ "${os}" = "opensuse" ]; then
        dist=`echo ${VERSION_ID}`
      elif [ "${os}" = "opensuse-leap" ]; then
        os=opensuse
        dist=`echo ${VERSION_ID}`
      else
        dist=`echo ${VERSION_ID} | awk -F '.' '{ print $1 }'`
      fi

    elif [ `which lsb_release 2>/dev/null` ]; then
      # get major version (e.g. '5' or '6')
      dist=`lsb_release -r | cut -f2 | awk -F '.' '{ print $1 }'`

      # get os (e.g. 'centos', 'redhatenterpriseserver', etc)
      os=`lsb_release -i | cut -f2 | awk '{ print tolower($1) }'`

    elif [ -e /etc/oracle-release ]; then
      dist=`cut -f5 --delimiter=' ' /etc/oracle-release | awk -F '.' '{ print $1 }'`
      os='ol'

    elif [ -e /etc/fedora-release ]; then
      dist=`cut -f3 --delimiter=' ' /etc/fedora-release`
      os='fedora'

    elif [ -e /etc/redhat-release ]; then
      os_hint=`cat /etc/redhat-release  | awk '{ print tolower($1) }'`
      if [ "${os_hint}" = "centos" ]; then
        dist=`cat /etc/redhat-release | awk '{ print $3 }' | awk -F '.' '{ print $1 }'`
        os='centos'
      elif [ "${os_hint}" = "scientific" ]; then
        dist=`cat /etc/redhat-release | awk '{ print $4 }' | awk -F '.' '{ print $1 }'`
        os='scientific'
      else
        dist=`cat /etc/redhat-release  | awk '{ print tolower($7) }' | cut -f1 --delimiter='.'`
        os='redhatenterpriseserver'
      fi

    else
      aws=`grep -q Amazon /etc/issue`
      if [ "$?" = "0" ]; then
        dist='6'
        os='aws'
      else
        unknown_os
      fi
    fi
  fi

  if [[ ( -z "${os}" ) || ( -z "${dist}" ) ]]; then
    unknown_os
  fi

  # remove whitespace from OS and dist name
  os="${os// /}"
  dist="${dist// /}"

  echo "Detected operating system as ${os}/${dist}."
}

finalize_yum_repo ()
{
  echo "Installing pygpgme to verify GPG signatures..."
  yum install -y pygpgme --disablerepo='gitlab_gitlab-ee'
  pypgpme_check=`rpm -qa | grep -qw pygpgme`
  if [ "$?" != "0" ]; then
    echo
    echo "WARNING: "
    echo "The pygpgme package could not be installed. This means GPG verification is not possible for any RPM installed on your system. "
    echo "To fix this, add a repository with pygpgme. Usualy, the EPEL repository for your system will have this. "
    echo "More information: https://fedoraproject.org/wiki/EPEL#How_can_I_use_these_extra_packages.3F"
    echo

    # set the repo_gpgcheck option to 0
    sed -i'' 's/repo_gpgcheck=1/repo_gpgcheck=0/' /etc/yum.repos.d/gitlab_gitlab-ee.repo
  fi

  echo "Installing yum-utils..."
  yum install -y yum-utils --disablerepo='gitlab_gitlab-ee'
  yum_utils_check=`rpm -qa | grep -qw yum-utils`
  if [ "$?" != "0" ]; then
    echo
    echo "WARNING: "
    echo "The yum-utils package could not be installed. This means you may not be able to install source RPMs or use other yum features."
    echo
  fi

  echo "Generating yum cache for gitlab_gitlab-ee..."
  yum -q makecache -y --disablerepo='*' --enablerepo='gitlab_gitlab-ee'
}

finalize_zypper_repo ()
{
  zypper --gpg-auto-import-keys refresh gitlab_gitlab-ee
  zypper --gpg-auto-import-keys refresh gitlab_gitlab-ee-source
}

main ()
{
  detect_os
  curl_check


  yum_repo_config_url="https://packages.gitlab.com/install/repositories/gitlab/gitlab-ee/config_file.repo?os=${os}&dist=${dist}&source=script"

  if [ "${os}" = "sles" ] || [ "${os}" = "opensuse" ]; then
    yum_repo_path=/etc/zypp/repos.d/gitlab_gitlab-ee.repo
  else
    yum_repo_path=/etc/yum.repos.d/gitlab_gitlab-ee.repo
  fi

  echo "Downloading repository file: ${yum_repo_config_url}"

  curl -sSf "${yum_repo_config_url}" > $yum_repo_path
  curl_exit_code=$?

  if [ "$curl_exit_code" = "22" ]; then
    echo
    echo
    echo -n "Unable to download repo config from: "
    echo "${yum_repo_config_url}"
    echo
    echo "This usually happens if your operating system is not supported by "
    echo "packagecloud.io, or this script's OS detection failed."
    echo
    echo "You can override the OS detection by setting os= and dist= prior to running this script."
    echo "You can find a list of supported OSes and distributions on our website: https://packages.gitlab.com/docs#os_distro_version"
    echo
    echo "For example, to force CentOS 6: os=el dist=6 ./script.sh"
    echo
    echo "If you are running a supported OS, please email support@packagecloud.io and report this."
    [ -e $yum_repo_path ] && rm $yum_repo_path
    exit 1
  elif [ "$curl_exit_code" = "35" -o "$curl_exit_code" = "60" ]; then
    echo
    echo "curl is unable to connect to packagecloud.io over TLS when running: "
    echo "    curl ${yum_repo_config_url}"
    echo
    echo "This is usually due to one of two things:"
    echo
    echo " 1.) Missing CA root certificates (make sure the ca-certificates package is installed)"
    echo " 2.) An old version of libssl. Try upgrading libssl on your system to a more recent version"
    echo
    echo "Contact support@packagecloud.io with information about your system for help."
    [ -e $yum_repo_path ] && rm $yum_repo_path
    exit 1
  elif [ "$curl_exit_code" -gt "0" ]; then
    echo
    echo "Unable to run: "
    echo "    curl ${yum_repo_config_url}"
    echo
    echo "Double check your curl installation and try again."
    [ -e $yum_repo_path ] && rm $yum_repo_path
    exit 1
  else
    echo "done."
  fi

  if [ "${os}" = "sles" ] || [ "${os}" = "opensuse" ]; then
    finalize_zypper_repo
  else
    finalize_yum_repo
  fi

  echo
  echo "The repository is setup! You can now install packages."
}

main


```