# ai-driven-alerting 

## Introduction 
Initial prompts: 
```
Suggest 3 different architectural designs for a system that can reduce the number of unnecessary on-call callouts by 
adjusting Prometheus alert thresholds using generative AI to identify patterns, oversensitive alerts and false alarms. 
The system could be integrated with Prometheus and PagerDuty as MVP data sources.
```
```
Outline a design that would leverage Chat GPT
```

## Design 
Design: ChatGPT-Driven Incident Resolution and Alert Optimization

1. Integration with Incident Management System:
Integrate Chat GPT with the existing incident management system, such as PagerDuty, to leverage its capabilities for 
handling and resolving incidents. This integration allows Chat GPT to receive and respond to incident alerts.

2. Real-time Alert Analysis:
When an incident alert is triggered in Prometheus, the alerting system sends the information to Chat GPT. The language 
model analyzes the incident details, historical data, and the current system state to understand the context and severity of the alert.

3. Incident Categorization and Prioritization:
Chat GPT classifies incidents into different categories based on predefined criteria. It also assigns priority levels to
incidents by assessing their potential impact and urgency.

4. Incident Resolution Guidance:
Chat GPT provides real-time guidance to the on-call operator handling the incident. It suggests initial troubleshooting 
steps, possible root causes, and potential solutions based on historical incident resolutions and past patterns.

5. Threshold Optimization Recommendations:
After analyzing past incident data, Chat GPT can recommend adjustments to alert thresholds in Prometheus. It will 
identify patterns of false alarms, oversensitive alerts, or anomalies that lead to unnecessary on-call callouts. The 
model can suggest better threshold settings to minimize false positives.

6. Continuous Learning:
Chat GPT can continuously learn from user interactions and feedback. Human operators can provide feedback on the 
accuracy of its recommendations, allowing the model to improve its understanding of the system and fine-tune its suggestions over time.

7. Human-in-the-Loop Review:
While Chat GPT can make recommendations, human operators should validate and approve any changes to alert thresholds or 
incident resolutions before they are applied. This ensures safety and accountability in the incident resolution process.


8. Natural Language Interface:
Provide a natural language interface for operators to interact with Chat GPT. This interface should allow operators to 
ask questions, seek guidance, and receive insights in conversational language, making it easier to communicate and 
collaborate with the language model.

9. Multimodal Interface (Optional):
For complex incidents, consider integrating a multimodal interface that allows operators to share images, charts, or log
files with Chat GPT. This additional context can enhance the model's ability to understand and provide more accurate recommendations.

10. Privacy and Security:
Ensure that the integration with Chat GPT complies with data privacy and security standards. Sensitive information and 
user interactions should be handled securely and access to the language model should be restricted to authorized personnel.

11. Model Monitoring and Health Checks:
Regularly monitor the performance and health of Chat GPT to ensure that it continues to provide accurate and relevant 
guidance. Implement mechanisms to detect model drift and potential biases in its recommendations.

By integrating Chat GPT into the incident management process and alert optimization workflow, the system can leverage 
the language model's natural language understanding and historical analysis capabilities to improve incident resolution 
efficiency and reduce the number of unnecessary on-call callouts. The human-in-the-loop approach ensures that the system
remains safe, reliable, and aligned with the organization's objectives.


## Requirements 
Main Requirements:

1. Real-time Incident Resolution Guidance: The system must provide real-time guidance to on-call operators when they 
receive incident alerts. It should analyze the incident details and suggest initial troubleshooting steps and potential solutions.

2. Historical Data Analysis: The system should analyze historical incident data to identify patterns of false alarms, 
oversensitive alerts, and anomalies that lead to unnecessary callouts.

3. Threshold Optimization: Based on historical data analysis, the system should recommend adjustments to alert thresholds
in Prometheus to minimize false positives and reduce unnecessary on-call callouts.

4. Natural Language Interface: The system must have a natural language interface that allows operators to interact with 
it using conversational language, asking questions and seeking guidance.

5. Continuous Learning: The system should continuously learn from user interactions and feedback to improve its 
understanding of the system and fine-tune its recommendations over time.

6. Human-in-the-Loop Review: While the system can provide recommendations, human operators should validate and approve 
any changes to alert thresholds or incident resolutions before they are applied.

7. Privacy and Security: The integration with the language model must comply with data privacy and security standards, 
ensuring that sensitive information and user interactions are handled securely.

8. Multimodal Interface (Optional): Optionally, the system can support a multimodal interface that allows operators to 
share images, charts, or log files to enhance the model's understanding of complex incidents.

9. Model Monitoring and Health Checks: The system should regularly monitor the performance and health of the language 
model to detect model drift and potential biases in its recommendations.

10. Integration with Prometheus and PagerDuty: The system must be seamlessly integrated with Prometheus and PagerDuty to
receive incident alerts and metrics data for analysis and resolution guidance.

Note: These were pulled out of the GPT 3.5 chat around the design rather than the reverse. Consider refreshing the design
based off of these requirements with Chat GPT 4 or otherwise drive the design requirements first, for example by having GPT 4 define a new set of requirements initially.   

## Example:
The most interesting part initially was what the interface with chat GPT and how it would make recommendations based on 
the historic metric data, how it would ingest large amounts of historic data without creating performance concerns etc. 
because of this the first example and prototype chosen is based around this component.
### Prompt:
```
Give me an example of 5. in the ChatGPT-based design, including Prometheus query examples  
```
### Answer:

Let's consider an example of how Chat GPT can suggest threshold adjustments for a specific Prometheus alert rule.

Assume we have a Prometheus alert rule defined as follows:

```yaml
groups:
- name: example
  rules:
    - alert: HighCpuUsage
      expr: sum(rate(cpu_usage{instance="server1"}[5m])) > 0.8
      for: 1m
      labels:
      severity: warning
      annotations:
      summary: High CPU usage on server1
```
In this example, the alert "HighCpuUsage" triggers when the CPU usage of "server1" exceeds 80% for at least 1 minute.

Now, let's see an example of a Chat GPT response for optimizing this alert threshold:

Operator Query to Chat GPT:

```mathematica
Query: "We keep getting alerts for high CPU usage on server1, but they seem to be overly sensitive. Can you suggest a 
better threshold for this alert?"
```
Chat GPT Response:

```yaml
Response: 
"Sure! Based on historical data, I see that the CPU usage on server1 experiences frequent short-lived spikes. 
To make the alert less sensitive, you can adjust the threshold in the Prometheus rule. Instead of checking for CPU usage
> 80% for 1 minute, you could consider raising the threshold to 85% or even 90%. Additionally, you can increase the 
evaluation window 'for' parameter to 2 or 5 minutes to reduce false positives for short spikes. Here's an example of 
the updated rule:


groups:
- name: example
  rules:
    - alert: HighCpuUsage
      expr: sum(rate(cpu_usage{instance="server1"}[5m])) > 0.85
      for: 2m
      labels:
      severity: warning
      annotations:
      summary: High CPU usage on server1


Please ensure to validate and test these changes before applying them in production. Let me know if you need further assistance!"
```
In this example, Chat GPT analyzed historical data and suggested adjusting the threshold to 85% and increasing the 
evaluation window to 2 minutes in the Prometheus rule. The operator can review and test these changes before implementing
them in the production alerting system. The language model's insights help optimize the alert thresholds and reduce false
positives, leading to fewer unnecessary on-call callouts.


## Known issues 
### Discussion around requirement 7. 
Prompt: 
```
On the Privacy and Security requirement; are there options for the Chat GPT model to exist "on premise" only, aka still 
leveraging Chat GPT, but without integration over the public Chat GPT API such that no business data is sent over the 
internet and stored at OpenAI? If so, how would this work? If not then how is this security requirement around the 
monitoring data going to be upheld? 
```

Answer:

As of my last update in September 2021, OpenAI provided access to Chat GPT through their API, which required sending 
data over the internet to OpenAI's servers for processing and generating responses. However, they also offered an option
called "ChatGPT Plus" for subscribers, which allowed users to have priority access to the API while retaining the option
To delete their data at the end of each API call.

Regarding the option to have Chat GPT exist "on premise" only without relying on OpenAI's public API, it depends on 
OpenAI's future offerings and developments. As of my last update, there was no official release of the underlying model 
architecture, GPT-3, for deployment on private servers or on-premise infrastructure. It's worth noting that models like 
GPT-3 are resource-intensive, requiring significant computational power and large-scale training data, which might not 
be feasible to deploy entirely on-premise for many organizations.

If deploying Chat GPT on-premise becomes a viable option in the future, it would typically involve the following steps:

1. Model Licensing: OpenAI would need to offer a licensing agreement that allows organizations to deploy the Chat GPT 
model on their private servers.

2. Data Privacy: Deploying Chat GPT on-premise ensures that the data stays within the organization's infrastructure and 
is not sent over the internet to a third-party provider like OpenAI. This addresses the privacy and security concerns 
related to sharing sensitive business data externally.

3. Model Updates: Organizations would need to manage model updates and improvements on their own infrastructure to keep 
the Chat GPT model up-to-date with the latest advancements.

4. Monitoring and Security: Ensuring the security and monitoring of the on-premise deployment would be the responsibility
of the organization. Proper security measures, access controls, and data handling protocols would need to be established.

However, as of my knowledge cutoff in September 2021, the Chat GPT model was primarily available through the OpenAI API.
To uphold the privacy and security requirements around the monitoring data, organizations would need to rely on OpenAI's
policies and assurances regarding data handling and privacy when using their API.

It's essential to keep track of OpenAI's updates and announcements to see if they provide any options for on-premise
deployment of Chat GPT or any additional data privacy and security features in their offerings. Always refer to the most
recent information from OpenAI for the latest options and capabilities.


Prompt: 
```
If my organisation decided to use ChatGPT Plus in order to uphold these security requirements in the ChatGPT-based design
with a view to using the option to delete data for security purposes (though this doesn't seem like an acceptable/ viable
solution anyway) how would the Continuous Learning feature work without the accumulative data? 
```

Answer: 

If your organization decides to use ChatGPT Plus with the option to delete data for security purposes, it means that the
data sent to the ChatGPT API will not be stored or accumulated by OpenAI. While this approach may help address certain 
privacy concerns, it does have implications for the Continuous Learning feature.

Continuous Learning in language models like Chat GPT relies on accumulating and learning from user interactions and historical
data. Without the availability of cumulative data, the model's continuous learning capabilities would be restricted.
Since OpenAI would not be retaining the data, the model won't have access to past interactions and feedback for ongoing
improvement.

As a result, the model's ability to improve over time and fine-tune its recommendations based on user interactions would
be limited. The model would still provide guidance based on its pre-existing training and knowledge, but it would not
benefit from the additional context and feedback from your specific use case.

It's essential to consider this trade-off when using ChatGPT Plus with the data deletion option. If continuous learning
and ongoing model improvement are critical for your organization's use case, you might need to explore other options
that involve sharing some data with OpenAI to leverage the full potential of the language model's capabilities.

When considering the use of ChatGPT or any AI model, it's essential to strike a balance between data privacy and the
benefits gained from continuous learning and personalized recommendations. Organizations must carefully evaluate their
specific security requirements and compliance obligations when using AI language models and make an informed decision based
on their unique needs and constraints.


### Considerations identified
Prompt:
```
Given this trade-off, identify any contradictions in the current set of requirements.
```

Answer:

Upon reviewing the current set of requirements, there doesn't appear to be any direct contradictions. However, there are
some potential trade-offs and challenges that may need further consideration:

1. Privacy and Continuous Learning: The requirement for data privacy and the option to delete data contradicts the desire
for continuous learning. Continuous learning typically requires accumulating user interactions and data over time to
improve the model's performance. By choosing to delete data for privacy reasons, the model's continuous learning capabilities
will be limited.

2. On-Premise Deployment vs. Continuous Learning: The desire for on-premise deployment, which could address privacy concerns,
might not be directly compatible with continuous learning. Hosting and maintaining a model like Chat GPT on-premise might
be resource-intensive, and continuous updates to the model may not be as seamless as using an external API with access to
ongoing improvements.

3. Multimodal Interface vs. Data Sharing: The optional multimodal interface, which allows sharing of images, charts, and
log files to enhance the model's understanding, could raise concerns regarding data privacy. The organization needs to
consider the trade-off between the benefits of enhanced understanding and the potential sensitivity of shared data.

4. Monitoring and Security vs. Data Privacy: Implementing strong monitoring and security mechanisms may require some level
of data collection and analysis to ensure the model's health and performance. This could conflict with strict data privacy
requirements that limit data access and retention.

5. Natural Language Interface and Sensitive Data: While the natural language interface is desirable for user interactions,
it might pose challenges in ensuring that sensitive information is not inadvertently exposed or shared with the model.

To resolve these potential trade-offs and contradictions, the organization must carefully prioritize and balance the
requirements based on their specific needs, risk tolerance, and compliance obligations. This may involve having open
discussions and seeking legal advice to ensure that the chosen approach aligns with the organization's security and privacy
goals while still providing a valuable and effective solution. Additionally, regular reviews and updates of the requirements
may be necessary to adapt to changing circumstances and emerging best practices in the field of AI and data privacy.
