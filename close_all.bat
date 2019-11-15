@echo off

taskkill /im DBServer.exe

taskkill /im ExternalGateway.exe

taskkill /im GameServer.exe

taskkill /im InnerGateway.exe

taskkill /im LoginServer.exe

taskkill /im simulate.exe

stop
