package trace

// Exporter indicates an available Span Exporter that is supported by this package.
type Exporter string

// ExporterOTLP is the OTLP exporter.
const ExporterOTLP Exporter = "OTLP"
