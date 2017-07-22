$module('basic', function () {
  var $static = this;
  $static.SomeGenerator = function () {
    var $current = 0;
    var $continue = function ($yield, $yieldin, $reject, $done) {
      while (true) {
        switch ($current) {
          case 0:
            $yield($t.fastbox(1, $g.________testlib.basictypes.Integer));
            $current = 1;
            return;

          case 1:
            $yield($t.fastbox(2, $g.________testlib.basictypes.Integer));
            $current = 2;
            return;

          case 2:
            $yield($t.fastbox(3, $g.________testlib.basictypes.Integer));
            $current = 3;
            return;

          default:
            $done();
            return;
        }
      }
    };
    return $generator.new($continue, false);
  };
  $static.TEST = $t.markpromising(function () {
    var $result;
    var $temp0;
    var $temp1;
    var counter;
    var entry;
    var s;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            s = $g.________testlib.basictypes.MapStream($g.________testlib.basictypes.Integer, $g.________testlib.basictypes.Integer)($g.basic.SomeGenerator(), function (s) {
              return $t.fastbox(s.$wrapped + 1, $g.________testlib.basictypes.Integer);
            });
            counter = $t.fastbox(0, $g.________testlib.basictypes.Integer);
            $current = 1;
            $continue($resolve, $reject);
            return;

          case 1:
            $temp1 = s;
            $current = 2;
            $continue($resolve, $reject);
            return;

          case 2:
            $promise.maybe($temp1.Next()).then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            entry = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 4;
              $continue($resolve, $reject);
              return;
            } else {
              $current = 5;
              $continue($resolve, $reject);
              return;
            }
            break;

          case 4:
            counter = $t.fastbox(counter.$wrapped + entry.$wrapped, $g.________testlib.basictypes.Integer);
            $current = 2;
            $continue($resolve, $reject);
            return;

          case 5:
            $resolve($t.fastbox(counter.$wrapped == 9, $g.________testlib.basictypes.Boolean));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
