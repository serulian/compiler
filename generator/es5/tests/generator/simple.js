$module('simple', function () {
  var $static = this;
  $static.SomeGenerator = function () {
    var $current = 0;
    var $continue = function ($yield, $yieldin, $reject, $done) {
      while (true) {
        switch ($current) {
          case 0:
            $yield($t.fastbox(false, $g.____testlib.basictypes.Boolean));
            $current = 1;
            return;

          case 1:
            $yield($t.fastbox(true, $g.____testlib.basictypes.Boolean));
            $current = 2;
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
            v = null;
            $current = 1;
            continue;

          case 1:
            $g.simple.SomeGenerator().then(function ($result0) {
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
            v = value;
            $current = 3;
            continue;

          case 6:
            $resolve(v);
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
