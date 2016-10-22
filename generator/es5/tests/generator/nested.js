$module('nested', function () {
  var $static = this;
  $static.AnotherGenerator = function () {
    var $current = 0;
    var $continue = function ($yield, $yieldin, $reject, $done) {
      while (true) {
        switch ($current) {
          case 0:
            $yield($t.fastbox(1, $g.____testlib.basictypes.Integer));
            $current = 1;
            return;

          case 1:
            if (true) {
              $current = 2;
              continue;
            } else {
              $current = 6;
              continue;
            }
            break;

          case 2:
            $yield($t.fastbox(2, $g.____testlib.basictypes.Integer));
            $current = 3;
            return;

          case 3:
            $current = 4;
            continue;

          case 4:
            $yield($t.fastbox(4, $g.____testlib.basictypes.Integer));
            $current = 5;
            return;

          case 6:
            $yield($t.fastbox(3, $g.____testlib.basictypes.Integer));
            $current = 7;
            return;

          case 7:
            $current = 4;
            continue;

          default:
            $done();
            return;
        }
      }
    };
    return $generator.new($continue);
  };
  $static.SomeGenerator = function () {
    var $result;
    var $current = 0;
    var $continue = function ($yield, $yieldin, $reject, $done) {
      while (true) {
        switch ($current) {
          case 0:
            $g.nested.AnotherGenerator().then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($yield, $yieldin, $reject, $done);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $yieldin($result);
            $current = 2;
            return;

          case 2:
            $yield($t.fastbox(5, $g.____testlib.basictypes.Integer));
            $current = 3;
            return;

          default:
            $done();
            return;
        }
      }
    };
    return $generator.new($continue);
  };
  $static.TEST = function () {
    var $result;
    var $temp0;
    var $temp1;
    var v;
    var value;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            v = $t.fastbox(0, $g.____testlib.basictypes.Integer);
            $current = 1;
            continue;

          case 1:
            $g.nested.SomeGenerator().then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $temp1 = $result;
            $current = 3;
            continue;

          case 3:
            $temp1.Next().then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 4;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 4:
            value = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 5;
              continue;
            } else {
              $current = 6;
              continue;
            }
            break;

          case 5:
            v = $t.fastbox(v.$wrapped + value.$wrapped, $g.____testlib.basictypes.Boolean);
            $current = 3;
            continue;

          case 6:
            $resolve($t.fastbox(v.$wrapped == 12, $g.____testlib.basictypes.Boolean));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
