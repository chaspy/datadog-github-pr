#!/bin/bash
# ref: https://docs.datadoghq.com/api/latest/metrics/#submit-metrics
api_key=${DATADOG_API_KEY}

export NOW="$(date +%s)"
curl -X POST "https://api.datadoghq.com/api/v1/series?api_key=${DATADOG_API_KEY}" \
-H "Content-Type: application/json" \
-d @- << EOF
{
  "series": [
    {
      "metric": "chaspy.tag.test",
      "points": [
        [
          "${NOW}",
          "1234.5"
        ]
      ],
      "tags": [
        "environment:test",
        "environment:staging",
        "author:chaspy"
      ]
    }
  ]
}
EOF

