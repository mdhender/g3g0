# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Shared run-time limit for the compat scripts (G3G0 #47). Source it from a
# run.sh:
#
#   . "${HERE}/../lib/timeout.sh"
#
# run_limited SECONDS ARGV0 BIN ARGS... runs BIN with argv[0] set to ARGV0
# and the given arguments, as the child of a small Perl supervisor. stdin,
# stdout and stderr are inherited unchanged.
#
# If BIN ends by itself, the supervisor exits with the status the shell would
# have seen for BIN: its exit code, or 128 + N if signal N killed it (134 for
# SIGABRT, 139 for SIGSEGV). The supervisor itself exits normally, so the
# shell prints no "Abort trap: 6" report of its own into the captured stderr.
#
# If BIN is still running after SECONDS, the supervisor sends it SIGTERM,
# then SIGKILL if it has not ended 2 seconds later, reaps it, and exits with
# TIMEOUT_STATUS. An alarm set before exec is not enough: the C programs die
# on SIGALRM, but the Go runtime ignores it, so a hung Go binary ran forever.
# Both die on SIGTERM, and nothing survives SIGKILL, so C and Go time out the
# same way and report the same status.
#
# TIMEOUT_STATUS is 124, as for GNU timeout. No G3 program and no part of the
# port exits 124 on its own, so in a status file it always means a timeout.

TIMEOUT_STATUS=124

run_limited() {
	TIMEOUT_STATUS="${TIMEOUT_STATUS}" perl -e '
		use strict;
		use POSIX ();
		my ($secs, $name, $bin, @args) = @ARGV;
		my $pid = fork;
		defined $pid or die "run_limited: fork: $!\n";
		if ($pid == 0) {
			{ exec { $bin } $name, @args; }
			print STDERR "run_limited: exec $bin: $!\n";
			POSIX::_exit(127);
		}
		my $done = eval {
			local $SIG{ALRM} = sub { die "timeout\n" };
			alarm $secs;
			waitpid($pid, 0);
			alarm 0;
			1;
		};
		if (!$done) {
			die $@ unless $@ eq "timeout\n";
			kill "TERM", $pid;
			for (1 .. 20) {
				exit $ENV{TIMEOUT_STATUS} if waitpid($pid, POSIX::WNOHANG()) == $pid;
				select(undef, undef, undef, 0.1);
			}
			kill "KILL", $pid;
			waitpid($pid, 0);
			exit $ENV{TIMEOUT_STATUS};
		}
		exit(($? & 127) ? 128 + ($? & 127) : $? >> 8);
	' "$@"
}
