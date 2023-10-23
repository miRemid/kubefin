// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
import React from "react";
import { Link } from "@cloudscape-design/components";
export function useDisclaimerFlashbarItem(onDismiss) {
  const email = "admin@kubefin.dev";
  const mailtoUrl = `mailto:${email}`;
  return {
    type: "info",
    dismissible: true,
    dismissLabel: "Dismiss message",
    onDismiss: () => onDismiss(),
    statusIconAriaLabel: "info",
    content: (
      <>
        Have questions or feedback on KubeFin? We're on{" "}
        <Link external href="https://kubefin.slack.com">
          Slack{" "}
        </Link>
        or email at{" "}
        <Link external href={mailtoUrl}>
          admin@kubefin.dev
        </Link>
      </>
    ),
  };
}
