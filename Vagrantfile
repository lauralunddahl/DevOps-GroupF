DO_BOX_URL = "https://github.com/devopsgroup-io/vagrant-digitalocean/raw/master/box/digital_ocean.box"

Vagrant.configure("2") do |config|
  config.vm.define "droplet1" do |droplet|
    droplet.vm.provider :digital_ocean do |provider, override|
      override.ssh.private_key_path = ENV["PRIVATE_KEY_PATH_DEVOPS"]
      override.vm.box = 'digital_ocean'
      override.vm.box_url = DO_BOX_URL
      override.nfs.functional = false
      override.vm.allowed_synced_folder_types = :rsync
      provider.token = ENV["DIGITAL_OCEAN_TOKEN"]
      provider.ssh_key_name = ENV["SSH_KEY_NAME"]
      provider.image = 'ubuntu-18-04-x64'
      provider.region = 'nyc1'
      provider.size = 's-1vcpu-1gb'
    end
  end
end