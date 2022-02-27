#!usr/bin/perl

if ($#ARGV+1 != 3) {
	print "\nUsage: build.pl cfg_path email_path bin_path\n";
	exit;
}

@COMMANDS = (
	'go mod tidy',
	"go build -o $ARGV[2]"
);

$ENV{'CFG_PATH'} = $ARGV[0];
$ENV{'TEMPLATE_PATH'} = $ARGV[1];
$e = $ARGV[2];

foreach $c (@COMMANDS) {
	system("$c");
}

print "Executing: $e\n";
system("tilix -e $e &");
