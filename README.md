# The Autobus Application

This is Autobus, a suite of applications who are destined to forever gather information about many GPSs.

# Running

## Attached Resources

If in doubt what a attached resource is, [read this](https://12factor.net/backing-services).

- [Nats](https://nats.io)
- [MongoDB](https://mongodb.org)

Installation of those is critical for running the application **without** `docker-compose`. If you do decide to use `docker-compose`, then there is no need to install any attached resource.

However, running Nats and MongoDB locally can give you greater control and inspection.

You're in for a treat!

First, install `go`, `gb`, `docker` and `docker-compose`. You really can't do anything without these.

Well, you _can_ leave out `docker` and `docker-compose`, but running a entire stack with one command is really appealing.

Next, clone this repository and, in the root folder, do:

```
$ GOOS=linux gb build
```

Next, you can bring the whole environment up with `docker-compose`. Just apply the correct .yml files, in the right order:

```
$ docker-compose -f docker-compose.yml -f development.yml build
$ docker-compose -f docker-compose.yml -f development.yml up -d
```

That's it!

Of course, each binary can be run independently. Build with `$ gb build`, and run each one of the applications accordingly. Take a look at the Environment section to see the different env vars needed to run each project. Also, take a look at the Architecture section to see how the different applications interact with one another. Or, if you like digging, check out the `docker-compose.yml` and `development.yml` files. Be sure to check the differences between the `production.yml` and `development.yml` files.

# Architecture

- The `autobus-core` application opens up a TCP server at port 9009 by default.
- When a GPS connects, it pushes data through the socket, which is then forwarded to the configured NATS, in the `gps.update` subject. The message is forwarded untouched.
  - The `autobus-platform` application forms a queue group under `queue.web.database`, on the subject `gps.update`.
  - The `autobus-platform` application can easily be scaled horizontally (see `AUTOBUS_PLATFORM_HORIZONTAL`) *and* vertically (start new ones, yay NATS!)
  - The `autobus-platform`, then, with the payload received from `autobus-core` via the `gps.update` subject, (tries to) parse and inserts the GPS update on the underlying MongoDB database. There's a capped (1kb, 500 documents) collection for transient data, and a cold collection for further storage.
- The `autobus-web` application, when requested, access the MongoDB database, querying the GPS messages table.
- The `autobus-web` application also creates bus stops through it's API.

Whew! Hope you now have a clue on what the applications main responsibility is.

# Environment

## Autobus Core

- `AUTOBUS_CORE_NATS_URL`: The NATS URL the Core will publish messages to. The subject name is `gps.update`.
- `AUTOBUS_CORE_DEBUG`: Enables debugging.
- `AUTOBUS_TCP_HOST`: Changes the TCP host. Default is `0.0.0.0:9009`
- `AUTOBUS_CORE_HANDLERS`: Tweaks client concurrency. The number of handlers is, in effect, the total numbers of connected clients the Hub can hold before buffering subsequent connections. You should increase this if the number of clients go higher. Default is 2048. (10k concurrent connections should be fine, given the server is able to handle that. Expect memory issues only with extreme concurrency/spikes)
- `AUTOBUS_CORE_ACCEPT`: Effectively, the number of concurrent goroutines accepting connections. Tweak this _should_ make clients be accepted faster; discretion is advised, though. Default is 1024.

## Autobus Platform
- `AUTOBUS_PLATFORM_HORIZONTAL`: Tweaks the number of goroutines that register callbacks on the NATS client. Tweaking this should parallellize the queue output rate, but this also increases the load on the database. Discretion is advised. Default is 1024.
- `AUTOBUS_PLATFORM_NATS_URL`: The NATS URL the platform will listen messages in. Should be the same subject name as Autobus Core. The queue group name is `queue.pgsql`.
- `AUTOBUS_PLATFORM_MONGO_URL`: The MongoDB servers it will insert GPS messages into. TODO: more details on the schema.

## Autobus Web
- `AUTOBUS_WEB_HOST`: Sets the host for the API **and** where the server will listen to incoming requests.
- `AUTOBUS_WEB_MONGO_URL`: Sets the MongoDB URL it will read from.

# Applications

## The Web API, schemas, that sort of thing

Note: represented as JSON for example values, each top level element represents a MongoDB collection.
This is the current version.

```
{
	stops: [
		{
			_id: ObjectID("hex"),
			name: "Santa Cruz",
			location: {
				type: "Point",
				coordinates: [150, 20]
			}
		}
	],

	gps_data_transient: [ /* gps data that is NOT long-lived */ ],
	gps_data: [ /* gps data that IS long lived */ ]
}
```

In the next implementations, we're aiming for this:

```
{
	// stops represent the bus stops along a path.
	// each line has n stops.
	stops: [
		{
			_id: "8abf716348cfd",
			name: "Santa Cruz",
			address: "RUA JOSE MARGARIDO COSTA",
			location: {
				type: "Point",
				coordinates: [-143.4183747, 22.88463]
			}
		}
	],

	// a line represents a entity that encompasses N buses.
	lines: [
		{
			_id: "81737471874",
			hours: ["08:00", "09:00"],
			stops: [
				{_id: "8abf716348cfd"}
			],
			route: {
				type: "MultiLineString",
				coordinates: [
					[-22, -44],
					[-22, -44],
					[-22, -44],
					[-22, -44],
					[-22, -44]
				]
			}
		}
	],

	gps_data_transient: [ /* gps data that is NOT long-lived */ ],
	gps_data: [ /* gps data that IS long lived */ ]
}
```
