// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     public/app/plugins/gen.go
// Using jennies:
//     TSTypesJenny
//     PluginTSTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.

import * as common from '@grafana/schema';

export const PanelCfgModelVersion = Object.freeze([0, 0]);

export interface PanelOptions {
  dedupStrategy: common.LogsDedupStrategy;
  enableLogDetails: boolean;
  prettifyLogMessage: boolean;
  showCommonLabels: boolean;
  showLabels: boolean;
  showTime: boolean;
  sortOrder: common.LogsSortOrder;
  wrapLogMessage: boolean;
}
