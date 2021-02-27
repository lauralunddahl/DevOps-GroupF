
Vagrant.configure("2") do |config|
    config.vm.box = 'digital_ocean'
    config.vm.box_url = "https://github.com/devopsgroup-io/vagrant-digitalocean/raw/master/box/digital_ocean.box"
    config.ssh.private_key_path = ENV["PRIVATE_KEY_PATH_DEVOPS"]

    config.vm.synced_folder "remote_files", "/vagrant", type: "rsync"

    config.vm.define "droplet1" do |droplet|
        droplet.vm.provider :digital_ocean do |provider, override|
	    provider.ssh_key_name = ENV["SSH_KEY_NAME"]
	    provider.token = ENV["DIGITAL_OCEAN_TOKEN"]
	    provider.image = 'ubuntu-18-04-x64'
	    provider.region = 'nyc1'
	    provider.size = 's-1vcpu-1gb'
      	    override.nfs.functional = false
    	end
    droplet.vm.hostname = "minitwit-server"
    droplet.vm.provision "shell", inline: <<-SHELL

    echo -e "\nVerifying that docker works ...\n"
    docker run --rm hello-world
    docker rmi hello-world

    echo -e "\nOpening port for minitwit ...\n"
    ufw allow 5000

    echo -e "\nOpening port for minitwit ...\n"
    echo ". $HOME/.bashrc" >> $HOME/.bash_profile

    echo -e "\nConfiguring credentials as environment variables...\n"
    echo "export DOCKER_USERNAME='***REMOVED***'" >> $HOME/.bash_profile
    echo "export DOCKER_PASSWORD='***REMOVED***'" >> $HOME/.bash_profile
    source $HOME/.bash_profile

    echo -e "\nVagrant setup done ..."
    echo -e "minitwit will later be accessible at http://$(hostname -I | awk '{print $1}'):5000"
    echo -e "The mysql database needs a minute to initialize, if the landing page is stack-trace ..."
    SHELL
    end
end