sysminerd
=========

sysminerd is a Go daemon that collects, transforms, and forwards Linux system metrics to other third party systems.  The input, transform, and output modules used are all configurable.  Initial support will focus on cpu, memory, and network input, with output support for [Graphite](https://github.com/graphite-project/).

# Configuration

The main daemon and modules are configured with yaml files.  

# Modules

## Input

Input modules are used for collecting system metrics.  At launch they will be initialized with their configuration details.  At the interval specified in the main configuration file the module will be queried for its list of metrics.  The input modules are responsible for setting the correct timestamps associated with the metrics.

TBD: As a design constraint we may require that input modules return metrics in a non-blocking manner.  So the input module would be responsible for maintaining an internal cache that can return its current metrics immediately. 

## Transform

A transform module can specify a list of input modules that it will mutate.  At launch they will be initialized with their configuration details. After the main daemon receives the list of metrics from the associated input modules they will be sent to the transform module.

## Output

Output modules are used for sending system metrics to other third party systems.  At launch they will be initialized with their configuration details.  After the list of metrics completes the transform stage, the list of metrics will be sent to the output modules.  If an output module cannot send the metrics it should send an error.  The main daemon will queue metrics based on the configuration setting specified.

# Monitoring

The daemon will provide an http API that can easily be queried.  It can return metrics on the daemon itself or metrics that it is collecting.  External monitoring tools can leverage this data also.
