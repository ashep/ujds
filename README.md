# DataPimp

**DataPimp** is a service which acts as a proxy between data providers and consumers. It accepts data items from 
producers and forwards them for processing to one or more consumers.

Every incoming data item has a kind which is described by a JSON schema. When **DataPimp** receives a data item, it 
checks whether the item is valid against corresponding schema and if it is, accepts the item, forwarding it to 
subscribed consumers.

Each consumer can accept or reject a data item, notifying **DataPimp** with additional information, such as a reason of
rejection and/or anything else. **DataPimp** stores each data item history. 

## HTTP API

## Changelog

## Authors

- [Oleksandr Shepetko](https://shepetko.com)
