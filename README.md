# `launchwith`

run CMD ARGS... with environment variables set from yaml CONFIG

	usage: launchwith [-expand] CONFIG CMD ARGS...
	  -expand : expand $SHELL_VARS in yaml string values

# motivation

Already made the app to read config exclusively from the env.
Other teams would like to use yaml files to configure.
Instead of pulling viper and trying to figure out precedance rules between environment (which viper doesn't really handle well), config files, flags, env-specific config files, just make that a separate step.

# prior work

I don't know, I banged this out in like an hour, there's probably something out there that already does this.

