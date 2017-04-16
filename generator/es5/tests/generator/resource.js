$module('resource', function () {
  var $static = this;
  this.$class('b8c8c08d', 'SomeResource', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.released = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
      return instance;
    };
    $instance.Release = function () {
      var $this = this;
      $this.released = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
      return;
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Release|2|fd8bc7c9<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.SomeGenerator = function (sr) {
    var $temp0;
    var $current = 0;
    var $resources = $t.resourcehandler();
    var $continue = function ($yield, $yieldin, $reject, $done) {
      $done = $resources.bind($done, false);
      while (true) {
        switch ($current) {
          case 0:
            $temp0 = sr;
            $resources.pushr($temp0, '$temp0');
            $yield($t.fastbox(2, $g.____testlib.basictypes.Integer));
            $current = 1;
            return;

          case 1:
            $resources.popr('$temp0');
            $yield($t.fastbox(40, $g.____testlib.basictypes.Integer));
            $current = 2;
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
    var i;
    var sr;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            sr = $g.resource.SomeResource.new();
            counter = $t.fastbox(0, $g.____testlib.basictypes.Integer);
            $current = 1;
            $continue($resolve, $reject);
            return;

          case 1:
            $temp1 = $g.resource.SomeGenerator(sr);
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
            i = $temp0.First;
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
            counter = $t.fastbox(counter.$wrapped + i.$wrapped, $g.____testlib.basictypes.Boolean);
            $current = 2;
            $continue($resolve, $reject);
            return;

          case 5:
            $resolve($t.fastbox(sr.released.$wrapped && (counter.$wrapped == 42), $g.____testlib.basictypes.Boolean));
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
