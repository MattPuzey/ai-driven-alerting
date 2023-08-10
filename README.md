# ai-driven-alerting

## Introduction
Initial prompt:

```
Suggest an architectural design for a system that can reduce the number of unnecessary on-call callouts by adjusting 
Prometheus alert thresholds using ChatGPT to identify patterns, oversensitive alerts and false alarms. The system could 
be integrated with Prometheus and/ or PagerDuty as MVP data sources.

Please also identify a list of main requirements to go with the design. 
```

## Architectural Design
1. Data Collection Layer
    * Prometheus Integration: Collects metrics and alert data, including the ones that triggered callouts.
    * PagerDuty Integration: Retrieves on-call schedules, incidents, and responses. It is used to correlate alerts with actions.
2.	Data Processing Layer
    * Preprocessing Module: Standardizes and cleans data from Prometheus and PagerDuty. May include normalization, missing value imputation, and categorization.
    * Feature Extraction Module: Extracts and engineers relevant features that can provide insights into patterns and correlations between alerts, metrics, and callouts.
3.	Analysis Layer
    * Pattern Recognition Module: Uses ChatGPT to identify common patterns, correlations, or sequences that are leading to false alarms or oversensitive alerts. Natural language processing can be employed here to help understand and categorize alerts.
    * Anomaly Detection Module: Applies machine learning algorithms to detect anomalies in the data which might signify a false alarm.
4.	Decision Support Layer
    * Threshold Adjustment Module: Based on insights from the Analysis Layer, automatically adjust alert thresholds in Prometheus or suggest adjustments to human operators.
    * Notification System: Integrated with PagerDuty to notify relevant team members of any detected patterns, false alarms, and suggested threshold adjustments.
5.	User Interface (UI) Layer
    * Dashboard: Allows users to view real-time statistics, identified patterns, and make manual adjustments to thresholds.
    * Report Generation Tool: Enables users to generate reports on the performance of the system, identified patterns, and any other relevant analytics.
6.	Security and Compliance Layer
    * Ensuring all data is handled in compliance with relevant legal and company regulations.
## Main Requirements
1.	Data Integration
    * Compatibility with Prometheus and PagerDuty APIs.
    * Real-time data synchronization.
2.	Data Quality and Preprocessing
    * Handling of missing or inconsistent data.
3.	Analysis Capabilities
    * Application of machine learning models like ChatGPT for pattern recognition.
    * Implementation of anomaly detection algorithms.
4.	Automation and Decision Support
    * Automatic adjustment of alert thresholds.
    * Alerting mechanisms for identified patterns or anomalies.
5.	User Interface
    * User-friendly dashboard with real-time monitoring capabilities.
    * Report generation with customizable templates.
6.	Security and Compliance
    * Data encryption and access control.
    * Compliance with relevant legal regulations, like GDPR.
7.	Scalability
    * Architecture should be designed to handle increasing volumes of data and alerts.
8.	Testing and Monitoring
    * Continuous monitoring of the system's performance.
    * Regular testing to ensure optimal performance and reliability.

By following this architectural design, the system would be able to learn from historical data and real-time patterns, recognize unnecessary alerts, and dynamically adjust thresholds, thereby reducing the number of unnecessary on-call callouts.




