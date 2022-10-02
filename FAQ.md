# FAQ

- [Why should I use AutoVPN?](#why-should-i-use-autovpn)
- [Why should I NOT use AutoVPN?](#why-should-i-not-use-autovpn)
- [I want to watch Netflix/Disney+ in another country, will this tool help me?](#i-want-to-watch-netflixdisney-in-another-country-will-this-tool-help-me)
- [Will this tool hide me from hackers?](#will-this-tool-hide-me-from-hackers)
- [Will this tool hide me from the government?](#will-this-tool-hide-me-from-the-government)

## Why should I use AutoVPN?

### It's cheaper

Instead of paying over $10 per month, you only pay for what you use. Meaning
you will only spend tens of cents per month for a slightly better service.

### No logs

The installation sets OpenVPN to not log anything on the server. If that's not
enough, the whole VPN server is automatically destroyed after you disconnect,
so it will leave little trace of your activities.

### Better privacy (if you know what you're doing)

You have better control over your VPN servers and better oversight over the tech
stacks behind them. You can choose VPS providers that use secure and updated
virtualization technology for maximum protection against hackers.

As of 2022-09-01, Linode is the recommended default, since they use KVM for
hardware virtualization, which is up-to-date and has fewer security
vulnerabilities compared to other vendors virtualization software.

That said: other vendors should be fine for single-sessions, despite
vulnerabilities.

## Why should I NOT use AutoVPN?

### Fewer countries to choose from

Most VPN providers have a lot of countries to connect to, almost all of them
in-fact. With AutoVPN, the countries and regions you can choose from is limited
by the locations of the provider's data centers, so the specific country you want
to connect to may not be available.

### Doesn't protect from state actors

This alone won't hide your activity from state actors.

### May not be able to connect to the VPS provider from some countries

Some countries may block the VPS resources needed and thus this may not work in
these countries.

## I want to watch Netflix/Disney+ in another country, will this tool help me?

Yes! This tool will let you connect to datacenters across the world and spoof
your IP address in the process. Those websites won't know you're browsing from a
different country.

## Will this tool hide me from hackers?

Partially yes. For most people this will be enough to hide your activity on public
Wi-Fi at e.g. cafés, airports, etc. If you're browsing from your home Wi-Fi, it
won't make a difference.

**Remember**: this won't help against other attack vectors, like rubber duckys,
phishing emails, password leaks etc. To mitigate that you must revise your
overall operational security (OpSec).

## Will this tool hide me from X government?

No. If they want to find you, they *will* find you.

Even if the server gets destroyed and even if they're not somehow tapping into your
traffic (which we can assume), they could probably access your billing information,
the ip addresses and the timestamps for the servers you create for your sessions.

Or break your fingers, that's also a possibility.