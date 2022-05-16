#!/bin/bash

find .|grep md$|xargs -n1 tocenize -max 6 -min 2
