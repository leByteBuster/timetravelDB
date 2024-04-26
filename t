[33mcommit dc022bb582047e91e86d6906b2825e691e40d8d4[m[33m ([m[1;36mHEAD[m[33m)[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Apr 18 11:50:58 2024 +0200

    Added TODO to API

[33mcommit 8d05461cd01b24d918aa3c90a327ab2199262dac[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Apr 10 14:19:21 2024 +0200

    Refactor and clean up neo4j data loader
    
    Functions are splittet into smaller parts into more expressive ones doing less stuff for more clarity.

[33mcommit 1aa97ea7fcecb32946883f99414324709a1373c5[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Apr 10 14:12:24 2024 +0200

    Clean up of the CLI commands
    
    More concise and consistent naming of commands.

[33mcommit d187de166e7a75f5c6213ca18647a5bc75e29a14[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Apr 4 23:09:22 2024 +0200

    Some useful comments for the data generator

[33mcommit 8e4ec16309622e470b54106b5a49a1604338dc1c[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Jul 10 02:10:25 2023 +0200

    Final commit
    
    I added the raw data of the benchmarks under benchmark/test-data.
    furthermore i added the raw results of the response time in a pdf under
    benchmark/results. The data of the integration tests are manipulated
    within the tests so the tests work again. Since datetimes are stored
    as datetime and not as string anymore, this had to be changed for the
    integration test data. Right now we do this through Neo4j write
    commands. Do a new backup in the future and exchange the backup data.

[33mcommit 2a42f3e65ae4e5a7f06097ec8ef7a3f54dcbb3eb[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Fri Jun 23 14:22:17 2023 +0200

    WIP: undo some if this stuff

[33mcommit d63d12468bd18704d38e5d382d818c0ec39d5950[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Fri Apr 28 10:43:41 2023 +0200

    Benchmark
    
    This commit introduces benchmarking on base of different data templates
    which generate different datasets. These datasets are then loaded both,
    in a only neo4j manner as well as in the ttdb manner to allow performing
    equivalent cypher and ttql queries for performance comparison reasons.
    The benchmark results are saved to .csv files. (Each query is run 100
    times after a warmup of 10 times)

[33mcommit 974297fa373c774eebc53ae96d39ca15c162946a[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Tue Apr 25 10:45:00 2023 +0200

    Restructuring. Now we have directories api, query-processor,
    database-engine. Also started benchmarking

[33mcommit a0d6c2bfc1186389a913b7edd64792baf7414ac0[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Apr 24 10:19:10 2023 +0200

    WIP2:serveral changes

[33mcommit f51cad062bac3dc1d339845b464c1ccef46b513c[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Apr 5 19:38:22 2023 +0200

    WIP: several changes

[33mcommit e472a098f1b78c08ce42131478bc5adaed5bbcb5[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Apr 3 19:26:51 2023 +0200

    Small test fix

[33mcommit 030b91cb271a053b7a0f035acfb1ee79dcbaa2ff[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Fri Mar 24 10:06:09 2023 +0100

    Cange functionality to allow multiple node labels, single edge label in template
    
    allow multiple node labels, single edge label (because neo4j only allows
    a single edge type)

[33mcommit f47dc44b325b15b24103eb93e69060952b22a0da[m[33m ([m[1;33mtag: [m[1;33mv0.9.3[m[33m)[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Mar 15 12:19:31 2023 +0100

    add some more functionality to CLI-APi

[33mcommit c2d78d2ad3adbd3dd09feee8c8bf13bc688955cb[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Mar 15 12:19:04 2023 +0100

    Add database check for timescaledb in setup_environment script

[33mcommit 54dc6aa762b86a3c2764f883f74cce8e1d3581ac[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Mar 15 12:16:26 2023 +0100

    Updating Data Loaders and Data Generator for correct Template usage
    
    Especially the only-neo4j data adapter had to be adapted to the
    generated Data. Data is now generated and loaded correctly by both
    TTDB and only-Neo4j.

[33mcommit 678747ef5d4d325b2beeb7cf69f0d027ece5627f[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Tue Mar 14 16:25:37 2023 +0100

    update data-generator, allow defining data template in a yaml file
    
    The template of the generated data can now be defined in a yaml
    file. The data generated will be in the structure described by
    the template. Note: the generated data contains a list of labels
    now instead of a single label field. We need to adapt this in the
    data adapters.

[33mcommit f3cfa0341e06f7f411ffc8fc47bec4779411a7a1[m[33m ([m[1;33mtag: [m[1;33mv0.9.2[m[33m, [m[1;33mtag: [m[1;33mlist[m[33m)[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Mar 13 10:51:28 2023 +0100

    Stabilize integration test
    
    The restore of the neo4j database seems to need some extra time
    which somehow cannot be ensured with checking for the readiness of the
    database via database state (via neo4j-admin database info)

[33mcommit 95cce4e2dfdb86ff043a258a06445a1fe6cf748c[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Fri Mar 10 17:41:44 2023 +0100

    Add support for conditoins on non-timeseries (non-binary) propertiers in neo4j
    
    Conditions on properties are only forwarded to timescale if they are
    actual time-series. We differ them with a prefix. all properties which
    start with ts_ or properties_ are expected to contain a time-series
    which is managed by timescaledb.

[33mcommit 6eb0266d412e202995a14deb2b646e2c376d2b65[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Fri Mar 10 16:11:01 2023 +0100

    Some renaming
    
    Some renaming to be more consistent

[33mcommit a2d46b4facf0e4766f5072c7ea4f3d96a19a75b1[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Fri Mar 10 13:13:47 2023 +0100

    Some minor changes, renaming

[33mcommit 6596a82260a4c565afee30625eae4b4d53a959c7[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Fri Mar 10 13:02:39 2023 +0100

    Prevent from trying to fetch time-series/properties if they arent
    
    If a property does not start with the prefix properties or ts then
    it is not a link to a time-series and nothing can be fetched. So
    in this case the timescale querying is skipped, the neo4j property
    is left untouched because it contains the complete data in this case.
    TODO: differ conditional cases differently on these properties

[33mcommit e290f06ecc048967473a1094decf345e7b4de692[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Fri Mar 10 12:43:48 2023 +0100

    Fixing time-series/property fetching
    
    fixed a false condition wich prevented time-series fetching

[33mcommit 300551bb14d8c6d451268a21ff39c260f050639c[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Fri Mar 10 12:22:10 2023 +0100

    Fix testing script, test main_test.go still not working
    
    It seems to be a problem to get main_test.go to work because
    of the neo4j restore which seems to make problems if run
    through a secondary call. Could not figure out why yet
    but seems to have something to do with the terminal
    environment.

[33mcommit 25a3e65c10397e054cdba37fb32023bfd1fa0c1a[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Mar 9 14:21:23 2023 +0100

    New generated test-data, merged helper files

[33mcommit c83b5298099a90adfcfa9d1ea49ed0ffad1c0980[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Mar 9 14:18:19 2023 +0100

    update and restructure datatabase-api. Unify them among packages
    
    All database access is now done via the databaseapi except for
    the only-neo4j data-adpter. this is still to be done. Some renaming
    and restructuring of the database-access has been done. Neo4j and
    Timescale database access is happening now over dedicated files in
    the databaseapi package. A config file config.yaml for the database-access has
    been introduced.

[33mcommit 12c41c2abd43fa166e15f30ee770283b77c4bd1d[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Mar 9 12:15:18 2023 +0100

    Huge cleanup and restructuring of the database-engine
    
    A lot of code duplication has been removed. Some renaming
    has been done to be more consistent. Still a lot of potential there.

[33mcommit 3faf6c1a8c35b0bc3817322ddcfc7c87341cbb9e[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Mar 9 12:09:52 2023 +0100

    Update scripts, try to merge logs in TestIntegration go test
    
    Update the scripts, nothing logical, just path stuff. Tried to
    merge the output of the tests triggered by TestIntegration.
    Unfortunately not successful yet.

[33mcommit 633899cc179b40ed0c5d79a486fcac9a55d9f51d[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Tue Mar 7 12:08:16 2023 +0100

    Extend the cleanup of the test-script, introduce script for
    test-environment setup
    
    Introducing a new script which only setups and starts the test
    envionment without runnign the tests for manual testing purposes.
    Extended the cleanup of the test-script. Also cleans up timescaledb
    before and after the script. just in case of leftovers.

[33mcommit 15265e5f87b5f888627b1e735540d26b7d0d5ac8[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Mar 6 21:40:42 2023 +0100

    Introducing case NON-SHALLOW RETURN and NON-SHALLOW WHERE; Refactoring
    
    Introducing the cases for NON-SHALLOW queries which have property
    lookups only in RETURN or only in WHERE.
    Quite some refactoring is done in this commit.
    Adding two integration test cases for the new introduced NON-SHALLOW query cases.

[33mcommit 48be32c844a12a3535b861fd4389365bda072ba6[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Mar 6 14:03:30 2023 +0100

    WIP: further cleanup, restructuring, documentation

[33mcommit 5ee67f85053bcdbbf19e9cd9a9ef21efd1f5ef3e[m[33m ([m[1;33mtag: [m[1;33mv0.9.0[m[33m)[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Tue Feb 28 19:00:32 2023 +0100

    WIP

[33mcommit 4459b4df7087da483965163a3be210a315bec7a3[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Feb 27 16:40:11 2023 +0100

    cleanup cache files from docker testing environment

[33mcommit 3951318c0e4150e774bc12c87692f72b1039f601[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Feb 27 16:27:10 2023 +0100

    fixup test automatisation

[33mcommit 49f14f789816a2c771941ee8690612d044f29153[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Feb 27 16:25:02 2023 +0100

    Add further functionality for non-shallow querying
    
    It is supported now to query non-shallow queries for
    queries without property lookups and for queries with
    property lookups in where and return clause.

[33mcommit d26b026efe8751e276e162a3ea226c893545a071[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Feb 27 16:23:27 2023 +0100

    add single-event json for examining puroposes of test data

[33mcommit 2984dffc71f72700cd9366fdf5637bd7dc9a6e6e[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Feb 27 16:22:43 2023 +0100

    update the testdata, remove the remaining ElementId fields

[33mcommit 44c88e399744ae9fb146e631ec44756966e65ee4[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Feb 27 16:19:04 2023 +0100

    Add script for automating the integration test
    
    The script is called by the go test TestIntegration.
    It then copies the backup files to the backup volumes
    of the docker containers, runs the neo4j backup
    (couldnt figure out how to do this in docker compose yet),
    runs docker compose, waits until the setup is ready and then
    runs the go tests for testing shallow and non-shallow queries.

[33mcommit 754fa132dc7f41abf07509b2f69169d3b5d90a08[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Sat Feb 25 15:19:07 2023 +0100

    docker-compose testing environment
    
    Setup a docker-compose file which with "docker-compose up" starts two testing containers, one
    for neo4j and one for timescaledb. Therefore it uses some backup data
    which is used for testing purposes. This is to be found in test-data.
    The dockerfile uses these files for restoring databases. Now the
    integration tests in api_test.go can be run by hand.

[33mcommit a11e08a6ee0c8ddf788e90eabe9af079c8f843e8[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Feb 23 10:50:06 2023 +0100

    Add backups of test data, adds field test data wannacry
    
    Adds test-data dictionary which contains a neo4j and timescale dump
    for testing purposes (the sql is postgres as well but use the .bak file)
    
    Also json files with wannacry event data is added.

[33mcommit 995041b714bead566646a6e1c18e6f9c3fb0ca1c[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Feb 22 10:06:08 2023 +0100

    Update TODOs

[33mcommit 569baaf5c6c74e52267a66b2dc7b26480defba2c[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Feb 22 10:00:22 2023 +0100

    Add functionality for case 2 (SHALLOW, lookup in WHERE & RETURN) and
    case 5 (NON SHALLOW, no lookup)
    
    Supports case 2 (SHALLOW, relevant lookup in WHERE & RETURN) and
    case 5 (NON SHALLOW, no relevant lookup) now. A relevant lookup
    - until now - is a lookup which includes a comparison. In the
    future there might be more possibilities for a relevant lookup.

[33mcommit 37a086395646890516185d671134dc56fa6c78ed[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Feb 22 09:54:38 2023 +0100

    Add a list of RETURN projections and a list of RETURN elements with no
    lookup to parse result
    
    Adds an list with all plain RETURN elements from the RETURN clause.
    For examle 'n','e' when "RETURN n, e, a.prop". Also adds a list with all
    projections which would include 'n','e','a.prop' in this case.

[33mcommit 6db127deeb1c8b4e8dada29acdbeaf1535ecfa21[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Feb 22 00:07:41 2023 +0100

    Some renaming and some small rearranging changes

[33mcommit 850b87fa92f60d66bcf79c958f29e8f341d33238[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Tue Feb 21 23:49:06 2023 +0100

    Adjusting options of the TTQL CLI, introducing hack for any()
    
    Makes the appearance of the TTQL CLI more mature. Also a hack
    is introduced so we can use "WHERE any(a.prop > 20)" because it is more
    intuitive. Under the hood it is still the same as writing
    "WHERE a.prop > 20" though. This should be changed in the future.

[33mcommit 73ef7c32e36f31cfe821af3979be9f077b175088[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Tue Feb 21 23:46:34 2023 +0100

    Add some more helpers to utils
    
    Adds a data structure with set functionality to utils as well
    as a function that pretty prints a map in the order of a passed
    key list

[33mcommit 0c59a1424b20fd1ba9a61fe54b8b1e8bc7d8793d[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Feb 20 17:08:29 2023 +0100

    Do some cleanup and some restructuring
    
    The existence check of a value for a property lookup comparison
    in a WHERE clause does not fetch data anymore but checks the existence
    in the timescale database and return just a boolean. This should be
    way faster.
    Furthermore the calculation of the LookupInfo for the lookup comparions
    inside WHERE clause are happening inside the Parser now when iterating
    over the insights so it does not have to happen twice. This still could
    be done in the listener which would be even better.

[33mcommit 6559cadf1d0a1baef28edb3e70c58bb166ae0e5f[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Feb 20 17:07:24 2023 +0100

    move helpers to utils package

[33mcommit 753e95ce2a27db7523e25451af3966eabc56c859[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Fri Feb 17 15:15:30 2023 +0100

    Implement SHALLOW secenario 3: lookup in WHERE, no lookup in RETURN
    
    Implementation of the case if a relevant property lookup is happening
    in the WHERE clause and only in the WHERE clause.

[33mcommit b2ae3fb9526b715356f1ca47a47c9694002d7485[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Feb 15 23:18:19 2023 +0100

    Clean up Parser and ttql_listener_impl
    
    merge all parse tree listeners into one. Some further cleanup
    inside the parser package.

[33mcommit ce3f837674a8c6a4f9fbe28b08167cc156c0845a[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Feb 15 20:08:40 2023 +0100

    Add little command line interface for querying TTQL queries

[33mcommit f21e7b64677f82d4c89f5cbdc4c5f2ff36df33f4[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Feb 15 20:08:07 2023 +0100

    Add some info files about neo4j objects and about visualization ideas

[33mcommit 8847e087f1e406ba8087000b6265e81aa8a43229[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Feb 15 20:04:25 2023 +0100

    Add some more information from the parser listeners to the ParseResult
    
    Add a map of relevant lookups in the WHERE clause and relevant lookups
    in the RETURN clause to the ParseResult. Relevant lookups are those
    which need an access to timescaleDB.

[33mcommit 35ad584dc57b311f93c49849a9852425f7a5c34b[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Feb 15 19:53:41 2023 +0100

    Add support for shallow TTQL queries w. & w.o. prop lookup in RETURN
    
    Add support for shallow TTQL queries with & without property lookup in RETURN clause.
    so Queries of the form:
      FROM ... TO ...  SHALLOW MATCH (a) RETURN *
      FROM ... TO ...  SHALLOW MATCH (a) RETURN a.someProperty
    are supported now

[33mcommit 6dd66ee2d6c6b837304e08d1876dd0ba31356551[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Feb 13 09:56:36 2023 +0100

    Git refactore parse tree listeners
    
    The PropertyClauseInsight Object for potential Property Lookups is saved
    as a map now. The keys are the ComparisonExpressionContext objects, so
    all its children are saved to it as a list via its
    PropertyClauseInsights.
    The ElementListener is introduced. It listenes to Variables of elements
    (nodes, edges) occouring in the query. It also checks if it is part
    of a MATCh, a WHERE or a RETURN clause.
    Some additional information ist listened to. For example, every
    comparison counts its PartialComparissonExpressions. Furthermore
    it is checked if there is a NullPrdicateExpression. It might be
    possible that this only works for one partial comparison expression
    yet. A reminder comment is put in place.

[33mcommit a33456e8172ef75286ec5979778173d561353234[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Feb 9 20:06:45 2023 +0100

    WIP

[33mcommit 95395cc4939b869c495f1ee95230eb767f5a7517[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Feb 9 20:06:06 2023 +0100

    Add listener MatchListener which holds the Match clause

[33mcommit 41614a2707220089b026f019a42f17d38e06a32a[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Feb 9 20:05:08 2023 +0100

    Extend test cases for testing valid queries

[33mcommit de490feda52821f3ab249b0a2892b1c12db71c68[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Feb 9 19:31:31 2023 +0100

    Some cleanup
    
    Throw out some manual testing, thorow out the old simplistic parsing
    through string splitting and other minor things

[33mcommit 279e78c663ad4a55e68c56e207ddc92b0dc59640[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Feb 9 19:12:03 2023 +0100

    Add testing for the parsing of some test valid and invalid TTQL queries
    
    Some valid TTQL queries are tested if they throw an error during the
    parsing process which they shouldnt. Some invalid TTQL queries are
    testd if they throw an error during parsing which they should.

[33mcommit 7b69b4a2b518caf61aae27ec643a6617aa901465[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Feb 9 18:58:12 2023 +0100

    WIP2: listeners implementation
    
    Added
    WhereListener which returns the Where Clause (only one for now)
    ReturnListener which returns the Return Clause (only one for now)
    timeListener is now TimeShallowListener and also returns if the query
    contains the SHALLOW Token

[33mcommit fd940b4117ccd950da13cfe88544ab7c56b7fa79[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Feb 8 20:25:34 2023 +0100

    WIP: listeners implementation
    
    Timelistener for the parse tree: returns FROM and TO DATETIME string
    PropertyExpressionListener: returns an object about all neccessary
    informtion (not working yet)
    PropertyLookuplistener: to be removed ?

[33mcommit c0af4ff35390943cc7b298c98c5ff605998bbd3f[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Tue Feb 7 15:22:15 2023 +0100

    update go.mod go.sum

[33mcommit fd4e287b5320b2913dbd97247df0299df004c83c[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Tue Feb 7 15:17:57 2023 +0100

    Add parser for timetravel Qurey Language (ttQL)
    
    This new subdirectory and package parsre contains the
    self defined TTQL.g4 file which defines the grammer of ttQL.
    It is base on OpenCypher. The auto generated files and interfaces
    of the the parser and lexer are put in a sub package ttql_listener.
    Another sub package "listeners" for listener implementations for the
    parser is added along with a first implementation of a listener
    for the PropertyLookup clause.

[33mcommit 2c9681847bddb728a0c5d9c48d51a13a4d84a4f4[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Tue Feb 7 15:15:54 2023 +0100

    move data-generator to own package

[33mcommit c946af0640138c6f1b7b4e3fe8f37cba6ffcbf68[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Tue Feb 7 15:11:04 2023 +0100

    create main subdirectory for main package

[33mcommit 036a0a46490e281fe120f8a488ac608903037214[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Tue Feb 7 15:09:24 2023 +0100

    move data-adapter-neo4j-only to own package

[33mcommit 5a8c05a3608ca84e7e7c551970523be9722a9250[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Tue Feb 7 14:50:17 2023 +0100

    move data adapter to own package

[33mcommit bf5dc09caa1653970c7f782d68f2351012b6f0a2[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Tue Feb 7 14:23:44 2023 +0100

    move api to own package

[33mcommit c453b38b64c0ebbd4508ddee61ab3b1c7563e568[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Feb 6 11:48:50 2023 +0100

    Merged modules into one

[33mcommit 26cedf89a33b29e6189ad4935af79d2285dd934a[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Feb 6 11:23:30 2023 +0100

    update gitignore, delete sub-modules

[33mcommit 82b2d5cdb4e6c2b5e73f654d66e019007cf6e97d[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Feb 6 10:10:27 2023 +0100

    WIP2: Add further functionality for api, start parsing
    
    Added further functionality to the api and started
    to implement parsing functionality through splitting
    the spring manually. This is expected to be replaced
    by the ANTLR implementation of the parser

[33mcommit 80adc2187d72fde54e08118e855d28122f4dbb41[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Jan 30 17:13:59 2023 +0100

    WIP: Include read functions for reading rows of Timescale for our api
    
    Added some functionality for reading single as well as multiple
    rows out of tables (timeseries) oftimescaledb. For this a Type
    TimeSeriesRow has been introduced. This is not tested yet, still
    WIP

[33mcommit 4d44ddf9fc47a04ab8ec4f4d3b133200fe5adfc1[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Jan 30 17:13:20 2023 +0100

    Fixed the data adatpter
    
    only timeseries data of edges were loaded

[33mcommit 07e63852395d1d8c20a9d7c189d9fe8f66bf0854[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Mon Jan 30 15:08:57 2023 +0100

    Cleaned data-adapter

[33mcommit 0dcde045ba45d8e027e2185c46cd0741d908cd2c[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Sat Jan 28 18:50:58 2023 +0100

    Multiple stuff
    
    - updating README
    - creating TODO
    - generated Files by file generator are loaded directly into the
      adapters
    - began to implement functions for the api to query neo4j AND tsDB

[33mcommit 9754e2faa212611738a8958641aa238e911acb88[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Jan 18 20:19:21 2023 +0100

    Move merged repository to subdirectory

[33mcommit ea3ee67ed7dc4324b1bb35469a58a5169ca4386c[m
Merge: 469aa98 f356972
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Jan 18 10:36:36 2023 +0100

    Merge remote-tracking branch 'temporal-graph-data-generator/main'
    
    Merging projects to one repository for better maintainance.

[33mcommit 469aa981bc79b371c4d869bdcd00f327bc7ef251[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Jan 18 10:32:18 2023 +0100

    Move merged repository to subdirectory

[33mcommit 8fe7abe66c7b97a5d87025ca3d919b905ac9f77b[m
Merge: 217d026 f60e6dc
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Jan 18 10:27:06 2023 +0100

    Merge remote-tracking branch 'api/main'
    
    Merge projects into one repository for maintainance

[33mcommit 217d026fe9f950dc98ceb66a191afd3173594145[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Jan 18 10:20:32 2023 +0100

    Move merged repository to subdirectory

[33mcommit 91bbebe12e7ccfbe7858e0dae0cd33410458f876[m
Merge: fa4e0a8 e8e3f3c
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Jan 18 10:13:21 2023 +0100

    Merge remote-tracking branch 'temporal-graph-data-neo4j-adapter/main'
    
    mergin projects to have one main repository

[33mcommit fa4e0a83b2a07c31f709bd813d85bb3b5812117b[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Jan 18 10:08:08 2023 +0100

    move merged repository to subdirectory

[33mcommit 26db7f6d457c519b563067de49b908e28d936d5d[m
Merge: 74d37ec e07ac11
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Jan 18 10:04:48 2023 +0100

    Merge remote-tracking branch 'timetraveldb-adapter/master'

[33mcommit 74d37ec45d399dd986426a22605457b563186187[m
Author: LexaTRex <40356781+LexaTRex@users.noreply.github.com>
Date:   Wed Jan 18 09:49:26 2023 +0100

    Initial commit

[33mcommit e07ac1129e9a85deacb3e657c23a1da64366fe2a[m[33m ([m[1;31mtimetraveldb-adapter/master[m[33m)[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Jan 18 09:26:20 2023 +0100

    Load data created by data generator into neo4j and timescaleDB.
    
    The Data is splitted into 1. the graph structure of the data along
    with the basic node/edge infos and property structure and
    2. the time-series like data structures which represent the values
    of the properties.
    Part 1. is loaded into neo4j, with uuids as property values. These
    uuids are used as table names for storing 2. into timescaleDB.

[33mcommit e8e3f3ce02ff6701ed324fc847b2bd22cc5a6999[m[33m ([m[1;31mtemporal-graph-data-neo4j-adapter/main[m[33m)[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Jan 12 12:20:58 2023 +0100

    Gain ability to read in edges and nodes created by the data_generator

[33mcommit f60e6dcc8d6dcdf8bef39cfa9d38e90391a85095[m[33m ([m[1;31mapi/main[m[33m)[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Sat Jan 14 21:35:34 2023 +0100

    Start of the API - not much progress now, mostly some thougths
    and structuring

[33mcommit fcb99db78b49d6f41e494ee3a7857a0900ce9629[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Sun Jan 15 18:02:33 2023 +0100

    initial commit

[33mcommit f35697212602460f316212dd1dbfa3ce4057bd8b[m[33m ([m[1;31mtemporal-graph-data-generator/main[m[33m)[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Sat Jan 14 21:29:21 2023 +0100

    Extend the data generator so that multiple values can be generated
    per property. Furthermore for nodes and edges time boundaries have
    to be passed. The restrictions for the time values of the property values
    are enforced now as well.

[33mcommit c2cab68465c129f79ef24f5e0babf81b59aa5c47[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Jan 12 12:47:22 2023 +0100

    Initial commit

[33mcommit 7dc0c506ef51bb6b555e5855b84acdf5eb6d0f28[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Thu Jan 5 21:14:42 2023 +0100

    Create neo4j nodes from json

[33mcommit a8f7e67d55f0f64117b21457c7186daaff281df1[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Jan 4 19:00:14 2023 +0100

    first commit

[33mcommit 96b2fc172b2b36770908f877715bfcf80aaf9f01[m
Author: Axel Schuster <axel_schuster_@gmx.net>
Date:   Wed Jan 4 18:02:06 2023 +0100

    initial commit
