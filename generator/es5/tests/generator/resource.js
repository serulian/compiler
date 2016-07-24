$module('resource', function () {
  var $static = this;
  this.$class('SomeResource', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve($t.box(false, $g.____testlib.basictypes.Boolean)).then(function (result) {
        instance.released = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.Release = function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $this.released = $t.box(true, $g.____testlib.basictypes.Boolean);
        $resolve();
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      return $t.createtypesig(['Release', 2, $g.____testlib.basictypes.Function($t.void).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.resource.SomeResource).$typeref()]);
    };
  });

  $static.SomeGenerator = function (sr) {
    var $temp0;
    var $current = 0;
    var $resources = $t.resourcehandler();
    var $continue = function ($yield, $yieldin, $reject, $done) {
      $done = $resources.bind($done);
      $reject = $resources.bind($reject);
      while (true) {
        switch ($current) {
          case 0:
            $temp0 = sr;
            $resources.pushr($temp0, '$temp0');
            $yield($t.box(2, $g.____testlib.basictypes.Integer));
            $current = 1;
            return;

          case 1:
            $resources.popr('$temp0').then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($yield, $yieldin, $reject, $done);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $result;
            $yield($t.box(40, $g.____testlib.basictypes.Integer));
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
    var $temp0;
    var $temp1;
    var counter;
    var i;
    var sr;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.resource.SomeResource.new().then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            sr = $result;
            counter = $t.box(0, $g.____testlib.basictypes.Integer);
            $current = 2;
            continue;

          case 2:
            $g.resource.SomeGenerator(sr).then(function ($result0) {
              $result = $result0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            $temp1 = $result;
            $current = 4;
            continue;

          case 4:
            $temp1.Next().then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 5;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 5:
            $result;
            i = $temp0.First;
            if ($t.unbox($temp0.Second)) {
              $current = 6;
              continue;
            } else {
              $current = 8;
              continue;
            }
            break;

          case 6:
            $g.____testlib.basictypes.Integer.$plus(counter, i).then(function ($result0) {
              counter = $result0;
              $result = counter;
              $current = 7;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 7:
            $result;
            $current = 4;
            continue;

          case 8:
            $promise.resolve($t.unbox(sr.released)).then(function ($result0) {
              return ($promise.shortcircuit($result0, true) || $g.____testlib.basictypes.Integer.$equals(counter, $t.box(42, $g.____testlib.basictypes.Integer))).then(function ($result1) {
                $result = $t.box($result0 && $t.unbox($result1), $g.____testlib.basictypes.Boolean);
                $current = 9;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 9:
            $resolve($result);
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
