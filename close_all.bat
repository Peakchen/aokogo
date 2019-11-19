@echo off

taskkill /f /im DBServer.exe

taskkill /f /im ExternalGateway.exe

taskkill /f /im GameServer.exe

taskkill /f /im InnerGateway.exe

taskkill /f /im LoginServer.exe

taskkill /f /im simulate.exe

stop
